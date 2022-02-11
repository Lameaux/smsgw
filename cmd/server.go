package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/routes"
)

var srv *http.Server //nolint:gochecknoglobals

const serverShutdownTimeout = 5 * time.Second

func startAPIServer(app *config.AppConfig) {
	srv = &http.Server{
		Addr:    ":" + app.Port,
		Handler: routes.Gin(app),
	}

	logger.Infow("starting server", "port", app.Port)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen: %s\n", err)
		}
	}()
}

func shutdownAPIServer() {
	logger.Infow("shutting down API server")

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("api server forced to shutdown: ", err)
	}

	logger.Infow("api server exiting")
}
