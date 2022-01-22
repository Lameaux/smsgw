package routes

import (
	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/handlers"
	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/providers/callbacks/sandbox"
	"euromoby.com/smsgw/internal/services"
	"github.com/gin-gonic/gin"
)

func Gin(app *config.AppConfig) *gin.Engine {
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

	cb := handlers.NewCallbackHandler(app)

	auth := middlewares.NewAuthenticator(app)

	m := r.Group("/v1/sms/messages")
	m.Use(auth.Authenticate)
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

	co := r.Group("/v1/sms/callbacks/outbound")
	co.Use(auth.Authenticate)
	{
		co.GET("", cb.ListCallbacks)
		co.POST("", cb.RegisterCallback)
		co.DELETE("/:id", cb.UnregisterCallback)
	}

	ci := r.Group("/v1/sms/callbacks/inbound")
	ci.Use(auth.Authenticate)
	{
		ci.GET("", cb.ListCallbacks)
		ci.GET("/:shortcode", cb.ListCallbacks)
		ci.POST("/:shortcode", cb.RegisterCallback)
		ci.PUT("/:shortcode", cb.UpdateCallback)
		ci.DELETE("/:shortcode", cb.UnregisterCallback)

	}

	ps := r.Group("/v1/sms/providers/sandbox")
	{
		sandbox.SetupRoutes(ps, inboundService, outboundService)
	}

	return r
}
