package main

import (
	coreworkers "github.com/Lameaux/core/workers"
	"github.com/Lameaux/smsgw/internal/billing"
	nm "github.com/Lameaux/smsgw/internal/notifications/models"
	"github.com/Lameaux/smsgw/internal/notifications/notifiers/http"
	deliveryworkers "github.com/Lameaux/smsgw/internal/notifications/workers/delivery"
	odp "github.com/Lameaux/smsgw/internal/outbound/processors/delivery"
	ows "github.com/Lameaux/smsgw/internal/outbound/workers/sending"

	"github.com/Lameaux/smsgw/internal/config"
	"github.com/Lameaux/smsgw/internal/outbound/connectors"
)

func makeWorkers(app *config.App) (workers []coreworkers.Worker) {
	c := connectors.NewConnectorRepository(app)
	w := ows.NewWorker(app, c, billing.NewStubBilling())
	workers = append(workers, w)

	n := http.NewNotifier(app)
	on := deliveryworkers.NewWorker(
		"OutboundDeliveryWorker",
		app,
		nm.DeliveryNotificationOutbound,
		odp.NewProcessor(n),
	)
	workers = append(workers, on)
	return
}
