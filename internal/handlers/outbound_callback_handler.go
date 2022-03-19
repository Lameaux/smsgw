package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/config"
)

type OutboundCallbackHandler struct {
	appConfig *config.AppConfig
}

func NewOutboundCallbackHandler(appConfig *config.AppConfig) *OutboundCallbackHandler {
	return &OutboundCallbackHandler{appConfig}
}

func (mc OutboundCallbackHandler) ListCallbacks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (mc OutboundCallbackHandler) RegisterCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (mc OutboundCallbackHandler) UnregisterCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
