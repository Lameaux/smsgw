package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/notifiers"
	"euromoby.com/smsgw/internal/providers/connectors"
	"euromoby.com/smsgw/internal/routes"
	"euromoby.com/smsgw/internal/workers"
	"golang.org/x/sync/errgroup"
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

func startWorkers(app *config.AppConfig) (*errgroup.Group, context.CancelFunc) {
	c := connectors.NewConnectorRepository(app)
	ow := workers.NewOutboundMessageWorker(app, c)

	n := notifiers.NewOutboundNotifier(app)
	on := workers.NewOutboundDeliveryWorker(app, n)

	ctx, cancel := context.WithCancel(context.Background())

	runners := []*workers.Runner{
		workers.NewRunner(ctx, ow),
		workers.NewRunner(ctx, on),
	}

	var g errgroup.Group

	for _, r := range runners {
		r := r
		g.Go(func() error {
			return r.Exec()
		})
	}

	return &g, cancel
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

func shutdownWorkers(g *errgroup.Group, cancel context.CancelFunc) {
	logger.Infow("shutting down workers")

	cancel()

	if err := g.Wait(); err != nil {
		logger.Errorw("error while stopping workers", "error", err)
	}
	logger.Infow("workers stopped")
}
