package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/notifiers"
	"euromoby.com/smsgw/internal/providers/connectors"
	"euromoby.com/smsgw/internal/routes"
	"euromoby.com/smsgw/internal/workers"
)

func main() {
	app := config.NewAppConfig()

	startWorkers(app)

	logger.Infow("starting server", "port", app.Port)
	routes.Gin(app).Run(":" + app.Port)
}

func startWorkers(app *config.AppConfig) {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	c := connectors.NewConnectorRepository(app)
	ow := workers.NewOutboundMessageWorker(app, c)

	n := notifiers.NewOutboundNotifier(app)
	on := workers.NewOutboundDeliveryWorker(app, n)

	ctx, cancel := context.WithCancel(context.Background())

	runners := []*workers.Runner{
		workers.NewRunner(ctx, ow),
		workers.NewRunner(ctx, on),
	}

	var wg sync.WaitGroup
	wg.Add(len(runners))

	for _, r := range runners {
		go func(r *workers.Runner) {
			defer wg.Done()

			r.Start()
		}(r)
	}

	go func(cancel context.CancelFunc) {
		<-sigChannel
		logger.Infow("the interrupt received, waiting for workers to stop")

		cancel()

		wg.Wait()
		logger.Infow("workers stopped")

		logger.Infow("shutting down")
		app.Shutdown()

		logger.Infow("exiting")
		os.Exit(0)
	}(cancel)
}
