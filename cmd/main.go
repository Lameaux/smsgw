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

const (
	workersCount = 2
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

	ow, err := workers.NewOutboundWorker(app, c)
	if err != nil {
		logger.Fatal(err)
	}

	var wg sync.WaitGroup

	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go func() {
			defer wg.Done()

			ow.Run()
		}()
	}

	go func() {
		<-sigChannel
		logger.Infow("The interrupt received. Waiting for workers to stop....")
		ow.Terminate()

		wg.Wait()
		logger.Infow("Workers stopped.")

		logger.Infow("Shutting down...")
		app.Shutdown()

		os.Exit(0)
	}()
}
