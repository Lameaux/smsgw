package handlers

import (
	"io"
	"net/http/httptest"

	"euromoby.com/smsgw/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type header struct {
	Key   string
	Value string
}

func performAnonRequest(f gin.HandlerFunc, method, path string, body io.Reader, headers ...header) *httptest.ResponseRecorder {
	r := gin.Default()
	r.Handle(method, path, f)

	return performRequest(r, method, path, body, headers...)
}

func performAuthRequest(f gin.HandlerFunc, method, path string, body io.Reader, headers ...header) *httptest.ResponseRecorder {
	r := authRouter()
	r.Handle(method, path, f)

	return performRequest(r, method, path, body, headers...)
}

func performRequest(r *gin.Engine, method, path string, body io.Reader, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func authRouter() *gin.Engine {
	router := gin.Default()
	a := middlewares.NewAuthenticator(TestAppConfig)
	router.Use(a.Authenticate)
	return router
}
