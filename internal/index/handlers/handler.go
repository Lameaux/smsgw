package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/config"
)

type Handler struct{}

type Response struct {
	AppName    string
	AppVersion string
	Health     string
}

func NewHandler() *Handler {
	return &Handler{}
}

func (i *Handler) Index(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		AppName:    config.AppName,
		AppVersion: config.AppVersion,
		Health:     "OK",
	})
}
