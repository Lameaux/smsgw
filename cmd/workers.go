package main

import (
	"context"
	"euromoby.com/smsgw/internal/billing"
	"euromoby.com/smsgw/internal/notifications"
	nm "euromoby.com/smsgw/internal/notifications/models"
	op "euromoby.com/smsgw/internal/outbound/processors"
	ow "euromoby.com/smsgw/internal/outbound/workers"

	"golang.org/x/sync/errgroup"

	"euromoby.com/core/logger"
	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/notifications/notifiers"
	"euromoby.com/smsgw/internal/providers/connectors"

	coreworkers "euromoby.com/core/workers"
)

var (
	workersGroup  errgroup.Group     //nolint:gochecknoglobals
	cancelWorkers context.CancelFunc //nolint:gochecknoglobals
)

func startWorkers(app *config.App) {
	ctx, cancel := context.WithCancel(context.Background())
	cancelWorkers = cancel

	runners := make([]*coreworkers.Runner, 0)

	c := connectors.NewConnectorRepository(app)
	w := ow.NewMessageWorker(app, c, billing.NewStubBilling())
	runners = append(runners, coreworkers.NewRunner(ctx, w))

	n := notifiers.NewDefaultNotifier(app)

	on := notifications.NewWorker(
		"OutboundDeliveryWorker",
		app,
		nm.DeliveryNotificationOutbound,
		op.NewDeliveryProcessor(n),
	)
	runners = append(runners, coreworkers.NewRunner(ctx, on))

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
