package main

import (
	"context"

	"golang.org/x/sync/errgroup"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/notifiers"
	"euromoby.com/smsgw/internal/processors"
	"euromoby.com/smsgw/internal/providers/connectors"
	"euromoby.com/smsgw/internal/workers"
)

var (
	workersGroup  errgroup.Group     //nolint:gochecknoglobals
	cancelWorkers context.CancelFunc //nolint:gochecknoglobals
)

func startWorkers(app *config.AppConfig) {
	ctx, cancel := context.WithCancel(context.Background())
	cancelWorkers = cancel

	runners := make([]*workers.Runner, 0)

	c := connectors.NewConnectorRepository(app)
	ow := workers.NewOutboundMessageWorker(app, c)
	runners = append(runners, workers.NewRunner(ctx, ow))

	n := notifiers.NewDefaultNotifier(app)

	on := workers.NewDeliveryNotificationWorker(
		"OutboundDeliveryWorker",
		app,
		models.MessageTypeOutbound,
		processors.NewOutboundDeliveryProcessor(n),
	)
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
