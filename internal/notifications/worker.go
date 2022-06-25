package notifications

import (
	"errors"
	"github.com/Lameaux/smsgw/internal/notifications/models"
	"time"

	"github.com/Lameaux/core/db"
	"github.com/Lameaux/core/logger"
	coremodels "github.com/Lameaux/core/models"
	corerepos "github.com/Lameaux/core/repos"
	"github.com/Lameaux/smsgw/internal/config"
	"github.com/Lameaux/smsgw/internal/notifications/notifiers"
)

const (
	notificationMaxAttempts = 5
)

type Worker struct {
	name        string
	app         *config.App
	messageType models.DeliveryNotificationType
	processor   Processor
}

func NewWorker(name string, app *config.App, messageType models.DeliveryNotificationType, processor Processor) *Worker {
	return &Worker{
		name:        name,
		app:         app,
		messageType: messageType,
		processor:   processor,
	}
}

func (w *Worker) Run() (bool, error) {
	tx, err := corerepos.Begin(w.app.Config.DBPool)
	if err != nil {
		return false, err
	}

	defer corerepos.Rollback(tx)

	notification, err := NewRepo(tx).FindOneForProcessing(w.messageType)
	if err != nil {
		if errors.Is(err, coremodels.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	logger.Infow("worker found new notification for processing", "worker", w.name, "delivery", notification)

	if err := w.processNotification(tx, notification); err != nil {
		return false, err
	}

	if err := corerepos.Commit(tx); err != nil {
		return false, err
	}

	return true, nil
}

func (w *Worker) processNotification(db db.Conn, n *models.DeliveryNotification) error {
	resp, err := w.processor.SendNotification(db, n)

	if err != nil {
		w.handleFailure(n, resp, err)
	} else {
		w.handleSuccess(n, resp)
	}

	return NewRepo(db).Update(n)
}

func (w *Worker) handleFailure(n *models.DeliveryNotification, resp *notifiers.NotifierResponse, err error) {
	if errors.Is(err, models.ErrSendFailed) || errors.Is(err, coremodels.ErrInvalidJSON) {
		n.LastResponse = resp.Body
	} else {
		errorText := err.Error()
		n.LastResponse = &errorText
	}

	if errors.Is(err, coremodels.ErrNotFound) || errors.Is(err, models.ErrMissingNotificationURL) {
		n.Status = models.DeliveryNotificationStatusFailed

		return
	}

	w.tryReschedule(n)
}

func (w *Worker) tryReschedule(n *models.DeliveryNotification) {
	if n.AttemptCounter >= w.MaxAttempts() {
		n.Status = models.DeliveryNotificationStatusFailed

		logger.Errorw("sending failed", "worker", w.name, "delivery", n)

		return
	}

	n.NextAttemptAt = coremodels.CalculateNextAttemptTime(n.AttemptCounter)
	n.AttemptCounter++
	logger.Infow("sending attempt failed, will try again later",
		"worker", w.name,
		"delivery", n,
		"next_attempt_at", n.NextAttemptAt,
		"attempt_counter", n.AttemptCounter,
	)
}

func (w *Worker) handleSuccess(n *models.DeliveryNotification, resp *notifiers.NotifierResponse) {
	n.Status = models.DeliveryNotificationStatusSent
	n.LastResponse = resp.Body

	logger.Infow("successfully sent", "worker", w.name, "delivery", n)
}

func (w *Worker) MaxAttempts() int {
	return notificationMaxAttempts
}

func (w *Worker) Name() string {
	return w.name
}

func (w *Worker) SleepTime() time.Duration {
	return w.app.WorkerSleep
}
