package main

import (
	"context"
	"github.com/Lameaux/smsgw/internal/billing"
	"github.com/Lameaux/smsgw/internal/notifications"
	nm "github.com/Lameaux/smsgw/internal/notifications/models"
	op "github.com/Lameaux/smsgw/internal/outbound/processors"
	ow "github.com/Lameaux/smsgw/internal/outbound/workers"

	"golang.org/x/sync/errgroup"

	"github.com/Lameaux/core/logger"
	"github.com/Lameaux/smsgw/internal/config"
	"github.com/Lameaux/smsgw/internal/notifications/notifiers"
	"github.com/Lameaux/smsgw/internal/providers/connectors"

	coreworkers "github.com/Lameaux/core/workers"
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
