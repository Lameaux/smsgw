package routes

import (
	coreconfig "github.com/Lameaux/core/config"
	"github.com/Lameaux/smsgw/internal/billing"
	"github.com/Lameaux/smsgw/internal/inbound"
	ohg "github.com/Lameaux/smsgw/internal/outbound/handlers/group"
	ohm "github.com/Lameaux/smsgw/internal/outbound/handlers/message"
	ohs "github.com/Lameaux/smsgw/internal/outbound/handlers/send"
	osg "github.com/Lameaux/smsgw/internal/outbound/services/group"
	osm "github.com/Lameaux/smsgw/internal/outbound/services/message"
	oss "github.com/Lameaux/smsgw/internal/outbound/services/send"
	"github.com/Lameaux/smsgw/internal/users"
	"github.com/gin-gonic/gin"

	coremiddlewares "github.com/Lameaux/core/middlewares"
	"github.com/Lameaux/smsgw/internal/config"
	ils "github.com/Lameaux/smsgw/internal/inbound/listeners/sandbox"
	ih "github.com/Lameaux/smsgw/internal/index/handlers"
	"github.com/Lameaux/smsgw/internal/middlewares"
	ols "github.com/Lameaux/smsgw/internal/outbound/listeners/sandbox"
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
		ils.Routes(ps, inboundService)
		ols.Routes(ps, outboundService)
	}

	return r
}
