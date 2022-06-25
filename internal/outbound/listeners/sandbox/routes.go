package sandbox

import (
	osm "github.com/Lameaux/smsgw/internal/outbound/services/message"
	"github.com/gin-gonic/gin"
)

func Routes(ps *gin.RouterGroup, outboundService *osm.Service) {
	handler := NewHandler(outboundService)

	ps.POST("/outbound/ack", handler.Ack)
}
