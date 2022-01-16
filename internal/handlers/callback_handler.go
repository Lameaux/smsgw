package handlers

import (
	"net/http"

	"euromoby.com/smsgw/internal/config"
	"github.com/gin-gonic/gin"
)

type CallbackHandler struct {
	appConfig *config.AppConfig
}

func NewCallbackHandler(appConfig *config.AppConfig) *CallbackHandler {
	return &CallbackHandler{appConfig}
}

func (mc CallbackHandler) ListCallbacks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (mc CallbackHandler) RegisterCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (mc CallbackHandler) UpdateCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (mc CallbackHandler) UnregisterCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
