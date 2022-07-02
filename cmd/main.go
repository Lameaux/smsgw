package main

import (
	"context"
	"github.com/Lameaux/core/executors"
	"github.com/Lameaux/core/httpserver"
	"github.com/Lameaux/smsgw/internal/routes"
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

	srv := httpserver.Start(&app.Config, routes.Gin(app))
	workers := makeWorkers(app)
	workersExecutor := executors.NewExecutor(workers)

	// Listen for the interrupt signal.
	<-ctx.Done()
	logger.Infow("the interrupt received, shutting down gracefully, press Ctrl+C again to force")
	stop()

	httpserver.Shutdown(srv, httpserver.ShutdownTimeout)
	workersExecutor.Shutdown()
	app.Config.Shutdown()

	logger.Infow("bye")
}
