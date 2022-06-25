package testhelpers

import (
	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/users"
	"io"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/middlewares"
)

type Header struct {
	Key   string
	Value string
}

func PerformAnonRequest(f gin.HandlerFunc, method, path string, body io.Reader, headers ...Header) *httptest.ResponseRecorder {
	r := gin.Default()
	r.Handle(method, path, f)

	return performRequest(r, method, path, body, headers...)
}

func PerformAuthRequest(app *config.App, f gin.HandlerFunc, method, path string, body io.Reader, headers ...Header) *httptest.ResponseRecorder {
	r := authRouter()
	r.Handle(method, path, f)

	return performRequest(r, method, path, body, headers...)
}

func performRequest(r *gin.Engine, method, path string, body io.Reader, headers ...Header) *httptest.ResponseRecorder {
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
	a := middlewares.NewAuthenticator(users.NewTestAuth())
	router.Use(a.Authenticate)

	return router
}
