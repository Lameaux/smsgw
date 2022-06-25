package main

import (
	"context"
	"github.com/Lameaux/smsgw/internal/billing"
	nm "github.com/Lameaux/smsgw/internal/notifications/models"
	"github.com/Lameaux/smsgw/internal/notifications/notifiers/http"
	deliveryworkers "github.com/Lameaux/smsgw/internal/notifications/workers/delivery"
	odp "github.com/Lameaux/smsgw/internal/outbound/processors/delivery"
	ows "github.com/Lameaux/smsgw/internal/outbound/workers/sending"

	"golang.org/x/sync/errgroup"

	"github.com/Lameaux/core/logger"
	"github.com/Lameaux/smsgw/internal/config"
	"github.com/Lameaux/smsgw/internal/outbound/connectors"

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
	w := ows.NewWorker(app, c, billing.NewStubBilling())
	runners = append(runners, coreworkers.NewRunner(ctx, w))

	n := http.NewNotifier(app)

	on := deliveryworkers.NewWorker(
		"OutboundDeliveryWorker",
		app,
		nm.DeliveryNotificationOutbound,
		odp.NewProcessor(n),
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
