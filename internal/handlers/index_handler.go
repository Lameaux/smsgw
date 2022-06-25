package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/config"
)

type IndexHandler struct{}

type IndexResponse struct {
	AppName    string
	AppVersion string
	Health     string
}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (i *IndexHandler) Index(c *gin.Context) {
	c.JSON(http.StatusOK, IndexResponse{
		AppName:    config.AppName,
		AppVersion: config.AppVersion,
		Health:     "OK",
	})
}
