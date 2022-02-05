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
	srv := startAPIServer(app)
	g, workersStop := startWorkers(app)

	// Listen for the interrupt signal.
	<-ctx.Done()
	logger.Infow("the interrupt received, shutting down gracefully, press Ctrl+C again to force")

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()

	shutdownAPIServer(srv)
	shutdownWorkers(g, workersStop)
	app.Shutdown()

	logger.Infow("bye")
}
