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

func shutdownWorkers(g *errgroup.Group, cancel context.CancelFunc) {
	logger.Infow("shutting down workers")

	cancel()

	if err := g.Wait(); err != nil {
		logger.Errorw("error while stopping workers", "error", err)
	}
	logger.Infow("workers stopped")
}
