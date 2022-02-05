package workers

import (
	"time"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/notifiers"
	"euromoby.com/smsgw/internal/repos"
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
	ctx, done := repos.DBTxContext()
	defer done()

	tx, err := w.app.DBPool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	notificationRepo := repos.NewDeliveryNotificationRepo(tx)
	notification, err := notificationRepo.FindOneForProcessing(models.MessageTypeOutbound)
	if err != nil {
		return false, err
	}

	if notification == nil {
		return false, nil
	}

	logger.Infow("worker found new notification for processing", "worker", w.Name(), "delivery", notification)

	message, messageOrder, err := w.findMessageWithOrder(tx, notification.MessageID)
	if err != nil {
		return false, err
	}

	if message == nil {
		logger.Errorw("message not found", "worker", w.Name(), "sms", message)
		notification.Status = models.DeliveryNotificationStatusFailed
	} else {
		w.sendToMerchant(notification, messageOrder, message)
	}

	if err = notificationRepo.Update(notification); err != nil {
		return false, err
	}

	if err = tx.Commit(ctx); err != nil {
		return false, nil
	}

	return true, nil
}

func (w *OutboundDeliveryWorker) findMessageWithOrder(tx db.Conn, id string) (*models.OutboundMessage, *models.MessageOrder, error) {
	outboundMessageRepo := repos.NewOutboundMessageRepo(tx)
	message, err := outboundMessageRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	if message == nil {
		return nil, nil, nil
	}

	messageOrderRepo := repos.NewMessageOrderRepo(tx)

	messageOrder, err := messageOrderRepo.FindByID(message.MessageOrderID)
	if err != nil {
		return nil, nil, err
	}

	if messageOrder == nil {
		return nil, nil, nil
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
	switch err {
	case models.ErrSendFailed, models.ErrInvalidJSON:
		n.LastResponse = resp.Body
	default:
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
	return 5
}

func (w *OutboundDeliveryWorker) SleepTime() time.Duration {
	return w.app.WorkerSleep
}
