package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"euromoby.com/core/logger"
	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/routes"
)

const serverShutdownTimeout = 5 * time.Second

func startAPIServer(app *config.App) *http.Server {
	srv := &http.Server{
		Addr:    ":" + app.Config.Port,
		Handler: routes.Gin(app),
	}

	logger.Infow("starting server", "port", app.Config.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen: %s\n", err)
		}
	}()

	return srv
}

func shutdownAPIServer(srv *http.Server) {
	logger.Infow("shutting down API server")

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("api server forced to shutdown: ", err)
	}

	logger.Infow("api server exiting")
}
