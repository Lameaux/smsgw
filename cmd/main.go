package main

import (
	"context"
	"os/signal"
	"syscall"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := config.NewAppConfig()
	logger.Infow("Starting", "app", app.AppName, "version", app.Version)

	startAPIServer(app)
	startWorkers(app)

	// Listen for the interrupt signal.
	<-ctx.Done()
	logger.Infow("the interrupt received, shutting down gracefully, press Ctrl+C again to force")
	stop()

	shutdownAPIServer()
	shutdownWorkers()
	app.Shutdown()

	logger.Infow("bye")
}
