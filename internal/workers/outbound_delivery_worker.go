package workers

import (
	"errors"
	"time"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/notifiers"
	"euromoby.com/smsgw/internal/repos"
)

const (
	outboundDeliveryMaxAttempts = 5
)

type OutboundDeliveryWorker struct {
	app      *config.AppConfig
	notifier *notifiers.OutboundNotifier
}

func NewOutboundDeliveryWorker(app *config.AppConfig, notifier *notifiers.OutboundNotifier) *OutboundDeliveryWorker {
	return &OutboundDeliveryWorker{
		app:      app,
		notifier: notifier,
	}
}

func (w *OutboundDeliveryWorker) Run() (bool, error) {
	tx, err := repos.Begin(w.app.DBPool)
	if err != nil {
		return false, err
	}

	defer repos.Rollback(tx)

	notificationRepo := repos.NewDeliveryNotificationRepo(tx)

	notification, err := notificationRepo.FindOneForProcessing(models.MessageTypeOutbound)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	logger.Infow("worker found new notification for processing", "worker", w.Name(), "delivery", notification)

	if err := w.processNotification(tx, notification); err != nil {
		return false, err
	}

	if err := repos.Commit(tx); err != nil {
		return false, err
	}

	return true, nil
}

func (w *OutboundDeliveryWorker) processNotification(tx db.Conn, notification *models.DeliveryNotification) error {
	notificationRepo := repos.NewDeliveryNotificationRepo(tx)

	message, messageOrder, err := w.findMessageWithOrder(tx, notification.MessageID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			logger.Errorw("message not found", "worker", w.Name(), "sms", message)

			notification.Status = models.DeliveryNotificationStatusFailed

			return notificationRepo.Update(notification)
		}

		return err
	}

	w.sendToMerchant(notification, messageOrder, message)

	return notificationRepo.Update(notification)
}

func (w *OutboundDeliveryWorker) findMessageWithOrder(tx db.Conn, id string) (*models.OutboundMessage, *models.MessageOrder, error) {
	outboundMessageRepo := repos.NewOutboundMessageRepo(tx)

	message, err := outboundMessageRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	messageOrderRepo := repos.NewMessageOrderRepo(tx)

	messageOrder, err := messageOrderRepo.FindByID(message.MessageOrderID)
	if err != nil {
		return nil, nil, err
	}

	return message, messageOrder, nil
}

func (w *OutboundDeliveryWorker) sendToMerchant(n *models.DeliveryNotification, messageOrder *models.MessageOrder, message *models.OutboundMessage) {
	resp, err := w.notifier.SendNotification(messageOrder, message)

	if err != nil {
		w.handleFailure(n, resp, err)
	} else {
		w.handleSuccess(n, resp)
	}
}

func (w *OutboundDeliveryWorker) handleFailure(n *models.DeliveryNotification, resp *notifiers.SendNotificationResponse, err error) {
	if errors.Is(err, models.ErrSendFailed) || errors.Is(err, models.ErrInvalidJSON) {
		n.LastResponse = resp.Body
	} else {
		errorText := err.Error()
		n.LastResponse = &errorText
	}

	w.tryReschedule(n)
}

func (w *OutboundDeliveryWorker) tryReschedule(n *models.DeliveryNotification) {
	if n.AttemptCounter >= w.MaxAttempts() {
		n.Status = models.DeliveryNotificationStatusFailed

		logger.Errorw("sending failed", "worker", w.Name(), "delivery", n)

		return
	}

	n.NextAttemptAt = models.CalculateNextAttemptTime(n.AttemptCounter)
	n.AttemptCounter++
	logger.Infow("sending attempt failed, will try again later",
		"worker", w.Name(),
		"delivery", n,
		"next_attempt_at", n.NextAttemptAt,
		"attempt_counter", n.AttemptCounter,
	)
}

func (w *OutboundDeliveryWorker) handleSuccess(n *models.DeliveryNotification, resp *notifiers.SendNotificationResponse) {
	n.Status = models.DeliveryNotificationStatusSent
	n.LastResponse = resp.Body

	logger.Infow("successfully sent", "worker", w.Name(), "delivery", n)
}

func (w *OutboundDeliveryWorker) Name() string {
	return "OutboundDeliveryWorker"
}

func (w *OutboundDeliveryWorker) MaxAttempts() int {
	return outboundDeliveryMaxAttempts
}

func (w *OutboundDeliveryWorker) SleepTime() time.Duration {
	return w.app.WorkerSleep
}
