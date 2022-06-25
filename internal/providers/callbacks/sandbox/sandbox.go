package sandbox

import (
	"github.com/Lameaux/smsgw/internal/inbound"
	osm "github.com/Lameaux/smsgw/internal/outbound/services/message"
	"github.com/gin-gonic/gin"
)

const (
	SandboxProviderID = "sandbox"
)

func SetupRoutes(ps *gin.RouterGroup, inboundService *inbound.Service, outboundService *osm.Service) {
	sandboxInbound := NewInboundHandler(inboundService)
	sandboxOutbound := NewOutboundHandler(outboundService)

	ps.POST("/inbound/message", sandboxInbound.ReceiveMessage)
	ps.POST("/outbound/ack", sandboxOutbound.Ack)
}
