package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Lameaux/core/logger"
	"github.com/Lameaux/smsgw/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := config.NewApp()
	logger.Infow("Starting", "app", config.AppName, "version", config.AppVersion)

	srv := startAPIServer(app)
	startWorkers(app)

	// Listen for the interrupt signal.
	<-ctx.Done()
	logger.Infow("the interrupt received, shutting down gracefully, press Ctrl+C again to force")
	stop()

	shutdownAPIServer(srv)
	shutdownWorkers()
	app.Config.Shutdown()

	logger.Infow("bye")
}
