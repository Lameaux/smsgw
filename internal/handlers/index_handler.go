package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/config"
)

type IndexHandler struct {
	appConfig *config.AppConfig
}

type IndexResponse struct {
	AppName string
	Version string
	Health  string
}

func NewIndexHandler(appConfig *config.AppConfig) *IndexHandler {
	return &IndexHandler{appConfig}
}

func (i *IndexHandler) Index(c *gin.Context) {
	c.JSON(http.StatusOK, IndexResponse{
		AppName: i.appConfig.AppName,
		Version: i.appConfig.Version,
		Health:  "OK",
	})
}
