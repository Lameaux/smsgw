package main

import (
	"context"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/notifiers"
	"euromoby.com/smsgw/internal/providers/connectors"
	"euromoby.com/smsgw/internal/workers"
	"golang.org/x/sync/errgroup"
)

var (
	workersGroup  errgroup.Group
	cancelWorkers context.CancelFunc
)

func startWorkers(app *config.AppConfig) {
	ctx, cancel := context.WithCancel(context.Background())
	cancelWorkers = cancel

	runners := make([]*workers.Runner, 0)

	c := connectors.NewConnectorRepository(app)
	ow := workers.NewOutboundMessageWorker(app, c)
	runners = append(runners, workers.NewRunner(ctx, ow))

	n := notifiers.NewOutboundNotifier(app)
	on := workers.NewOutboundDeliveryWorker(app, n)
	runners = append(runners, workers.NewRunner(ctx, on))

	for _, r := range runners {
		r := r
		workersGroup.Go(func() error {
			return r.Exec()
		})
	}
}

func shutdownWorkers() {
	logger.Infow("shutting down workers")

	cancelWorkers()

	if err := workersGroup.Wait(); err != nil {
		logger.Errorw("error while stopping workers", "error", err)
	}
	logger.Infow("workers stopped")
}
