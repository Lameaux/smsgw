package handlers

import (
	"net/http"

	"euromoby.com/smsgw/internal/config"
	"github.com/gin-gonic/gin"
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
