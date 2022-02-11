package sandbox

import (
	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/services"
)

const (
	SandboxProviderID = "sandbox"
)

func SetupRoutes(ps *gin.RouterGroup, inboundService *services.InboundService, outboundService *services.OutboundService) {
	sandboxInbound := NewInboundHandler(inboundService)
	sandboxOutbound := NewOutboundHandler(outboundService)

	ps.POST("/inbound/message", sandboxInbound.ReceiveMessage)
	ps.POST("/outbound/ack", sandboxOutbound.Ack)
}
