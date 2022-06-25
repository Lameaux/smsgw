package sandbox

import (
	"github.com/Lameaux/smsgw/internal/inbound"
	"github.com/gin-gonic/gin"
)

func Routes(ps *gin.RouterGroup, inboundService *inbound.Service) {
	handler := NewHandler(inboundService)

	ps.POST("/inbound/message", handler.ReceiveMessage)
}
