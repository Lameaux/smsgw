package workers

import (
	"errors"
	"time"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/notifiers"
	"euromoby.com/smsgw/internal/processors"
	"euromoby.com/smsgw/internal/repos"
)

const (
	notificationMaxAttempts = 5
)

type DeliveryNotificationWorker struct {
	name        string
	app         *config.AppConfig
	messageType models.MessageType
	processor   processors.DeliveryNotificationProcessor
}

func NewDeliveryNotificationWorker(name string, app *config.AppConfig, messageType models.MessageType, processor processors.DeliveryNotificationProcessor) *DeliveryNotificationWorker {
	return &DeliveryNotificationWorker{
		name:        name,
		app:         app,
		messageType: messageType,
		processor:   processor,
	}
}

func (w *DeliveryNotificationWorker) Run() (bool, error) {
	tx, err := repos.Begin(w.app.DBPool)
	if err != nil {
		return false, err
	}

	defer repos.Rollback(tx)

	notificationRepo := repos.NewDeliveryNotificationRepo(tx)

	notification, err := notificationRepo.FindOneForProcessing(w.messageType)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	logger.Infow("worker found new notification for processing", "worker", w.name, "delivery", notification)

	if err := w.processNotification(tx, notification); err != nil {
		return false, err
	}

	if err := repos.Commit(tx); err != nil {
		return false, err
	}

	return true, nil
}

func (w *DeliveryNotificationWorker) processNotification(db db.Conn, n *models.DeliveryNotification) error {
	resp, err := w.processor.SendNotification(db, n)

	if err != nil {
		w.handleFailure(n, resp, err)
	} else {
		w.handleSuccess(n, resp)
	}

	return repos.NewDeliveryNotificationRepo(db).Update(n)
}

func (w *DeliveryNotificationWorker) handleFailure(n *models.DeliveryNotification, resp *notifiers.NotifierResponse, err error) {
	if errors.Is(err, models.ErrSendFailed) || errors.Is(err, models.ErrInvalidJSON) {
		n.LastResponse = resp.Body
	} else {
		errorText := err.Error()
		n.LastResponse = &errorText
	}

	if errors.Is(err, models.ErrNotFound) || errors.Is(err, models.ErrMissingNotificationURL) {
		n.Status = models.DeliveryNotificationStatusFailed

		return
	}

	w.tryReschedule(n)
}

func (w *DeliveryNotificationWorker) tryReschedule(n *models.DeliveryNotification) {
	if n.AttemptCounter >= w.MaxAttempts() {
		n.Status = models.DeliveryNotificationStatusFailed

		logger.Errorw("sending failed", "worker", w.name, "delivery", n)

		return
	}

	n.NextAttemptAt = models.CalculateNextAttemptTime(n.AttemptCounter)
	n.AttemptCounter++
	logger.Infow("sending attempt failed, will try again later",
		"worker", w.name,
		"delivery", n,
		"next_attempt_at", n.NextAttemptAt,
		"attempt_counter", n.AttemptCounter,
	)
}

func (w *DeliveryNotificationWorker) handleSuccess(n *models.DeliveryNotification, resp *notifiers.NotifierResponse) {
	n.Status = models.DeliveryNotificationStatusSent
	n.LastResponse = resp.Body

	logger.Infow("successfully sent", "worker", w.name, "delivery", n)
}

func (w *DeliveryNotificationWorker) MaxAttempts() int {
	return notificationMaxAttempts
}

func (w *DeliveryNotificationWorker) Name() string {
	return w.name
}

func (w *DeliveryNotificationWorker) SleepTime() time.Duration {
	return w.app.WorkerSleep
}
