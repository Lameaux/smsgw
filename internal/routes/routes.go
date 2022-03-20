package routes

import (
	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/handlers"
	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/providers/callbacks/sandbox"
	"euromoby.com/smsgw/internal/services"
)

func Gin(app *config.AppConfig) *gin.Engine { //nolint:funlen
	r := gin.Default()
	r.Use(middlewares.Timeout(app.WaitTimeout))

	i := handlers.NewIndexHandler(app)

	r.GET("/", i.Index)
	r.GET("/health", i.Index)

	messageOrderService := services.NewMessageOrderService(app)
	send := handlers.NewSendHandler(messageOrderService)
	status := handlers.NewStatusHandler(messageOrderService)

	outboundService := services.NewOutboundService(app)
	outbound := handlers.NewOutboundHandler(outboundService)

	inboundService := services.NewInboundService(app)
	inbound := handlers.NewInboundHandler(inboundService)

	icb := handlers.NewInboundCallbackHandler(app)
	ocb := handlers.NewOutboundCallbackHandler(app)

	auth := middlewares.NewAuthenticator(app)

	v1sms := r.Group("/v1/sms")

	v1smsauth := r.Group("/v1/sms")
	v1smsauth.Use(auth.Authenticate)

	m := v1smsauth.Group("/messages")
	{
		m.POST("", send.SendMessage)

		m.GET("/status/search", status.Search)
		m.GET("/status/:id", status.Get)

		m.GET("/outbound/search", outbound.Search)
		m.GET("/outbound/:id", outbound.Get)

		m.GET("/inbound/:shortcode/search", inbound.Search)
		m.GET("/inbound/:shortcode/:id", inbound.Get)
		m.PUT("/inbound/:shortcode/:id/ack", inbound.Ack)
	}

	co := v1smsauth.Group("/callbacks/outbound")
	{
		co.GET("", ocb.GetCallback)
		co.POST("", ocb.RegisterCallback)
		co.PUT("", ocb.UpdateCallback)
		co.DELETE("", ocb.UnregisterCallback)
	}

	ci := v1smsauth.Group("/callbacks/inbound/:shortcode")
	{
		ci.GET("", icb.GetCallback)
		ci.POST("", icb.RegisterCallback)
		ci.PUT("", icb.RegisterCallback)
		ci.DELETE("", icb.UnregisterCallback)
	}

	ps := v1sms.Group("/providers/sandbox")
	{
		sandbox.SetupRoutes(ps, inboundService, outboundService)
	}

	return r
}
