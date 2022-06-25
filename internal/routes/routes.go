package routes

import (
	coreconfig "euromoby.com/core/config"
	"euromoby.com/smsgw/internal/billing"
	"euromoby.com/smsgw/internal/inbound"
	ohg "euromoby.com/smsgw/internal/outbound/handlers/group"
	ohm "euromoby.com/smsgw/internal/outbound/handlers/message"
	ohs "euromoby.com/smsgw/internal/outbound/handlers/send"
	osg "euromoby.com/smsgw/internal/outbound/services/group"
	osm "euromoby.com/smsgw/internal/outbound/services/message"
	oss "euromoby.com/smsgw/internal/outbound/services/send"
	"euromoby.com/smsgw/internal/users"
	"github.com/gin-gonic/gin"

	coremiddlewares "euromoby.com/core/middlewares"
	"euromoby.com/smsgw/internal/config"
	ih "euromoby.com/smsgw/internal/index/handlers"
	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/providers/callbacks/sandbox"
)

func Gin(app *config.App) *gin.Engine { //nolint:funlen
	r := gin.Default()

	switch app.Config.Env {
	case coreconfig.EnvTest:
		gin.SetMode(gin.TestMode)
	case coreconfig.EnvProduction:
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(coremiddlewares.Timeout(app.Config.WaitTimeout))

	i := ih.NewHandler()

	r.GET("/", i.Index)
	r.GET("/health", i.Index)

	sendHandler := ohs.NewHandler(oss.NewService(app, billing.NewStubBilling()))
	groupHandler := ohg.NewHandler(osg.NewService(app))

	outboundService := osm.NewService(app)
	outboundHandler := ohm.NewHandler(outboundService)

	a := users.NewStubAuth()

	inboundService := inbound.NewService(app, a)
	inboundHandler := inbound.NewHandler(inboundService)

	v1sms := r.Group("/v1/sms")

	authenticator := middlewares.NewAuthenticator(a)
	v1smsauth := r.Group("/v1/sms")
	v1smsauth.Use(authenticator.Authenticate)

	m := v1smsauth.Group("/messages")
	{
		m.POST("", sendHandler.SendMessage)

		m.GET("/group/search", groupHandler.Search)
		m.GET("/group/:id", groupHandler.Get)

		m.GET("/outbound/search", outboundHandler.Search)
		m.GET("/outbound/:id", outboundHandler.Get)

		m.GET("/inbound/search", inboundHandler.Search)
		m.GET("/inbound/:id", inboundHandler.Get)
		m.PUT("/inbound/:id/ack", inboundHandler.Ack)
	}

	ps := v1sms.Group("/providers/sandbox")
	{
		sandbox.SetupRoutes(ps, inboundService, outboundService)
	}

	return r
}
