package main

import (
	"context"
	"net/http"
	"time"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/routes"
)

func startAPIServer(app *config.AppConfig) *http.Server {
	srv := &http.Server{
		Addr:    ":" + app.Port,
		Handler: routes.Gin(app),
	}

	logger.Infow("starting server", "port", app.Port)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: %s\n", err)
		}
	}()

	return srv
}

func shutdownAPIServer(srv *http.Server) {
	logger.Infow("shutting down API server")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("api server forced to shutdown: ", err)
	}
	logger.Infow("api server exiting")
}
