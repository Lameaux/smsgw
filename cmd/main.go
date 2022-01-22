package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/providers/connectors"
	"euromoby.com/smsgw/internal/routes"
	"euromoby.com/smsgw/internal/workers"
)

func main() {
	app := config.NewAppConfig()
	connectorRepo := connectors.NewConnectorRepository()

	startWorkers(app, connectorRepo)

	logger.Infow("Starting server", "port", app.Port)
	routes.Gin(app).Run(":" + app.Port)
}

func startWorkers(app *config.AppConfig, c *connectors.ConnectorRepository) {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	ow := workers.NewOutboundMessageWorker(app, c)

	runners := []*workers.Runner{
		workers.NewRunner(ow),
		// workers.NewRunner(ow),
		// workers.NewRunner(ow),
	}

	var wg sync.WaitGroup
	wg.Add(len(runners))

	for _, r := range runners {
		go func(r *workers.Runner) {
			defer wg.Done()

			r.Start()
		}(r)
	}

	go func() {
		<-sigChannel
		logger.Infow("The interrupt received. Waiting for workers to stop....")

		for _, r := range runners {
			r.Stop()
		}

		wg.Wait()
		logger.Infow("Workers stopped.")

		logger.Infow("Shutting down...")
		app.Shutdown()

		logger.Infow("Exiting...")
		os.Exit(0)
	}()
}
