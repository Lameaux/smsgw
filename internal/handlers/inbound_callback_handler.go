package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/config"
)

type InboundCallbackHandler struct {
	appConfig *config.AppConfig
}

func NewInboundCallbackHandler(appConfig *config.AppConfig) *InboundCallbackHandler {
	return &InboundCallbackHandler{appConfig}
}

func (mc InboundCallbackHandler) GetCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (mc InboundCallbackHandler) RegisterCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (mc InboundCallbackHandler) UpdateCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (mc InboundCallbackHandler) UnregisterCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
