package workers

import (
	"errors"
	"euromoby.com/smsgw/internal/billing"
	"euromoby.com/smsgw/internal/notifications"
	"time"

	"euromoby.com/core/db"
	"euromoby.com/core/logger"
	"euromoby.com/smsgw/internal/config"
	nm "euromoby.com/smsgw/internal/notifications/models"
	"euromoby.com/smsgw/internal/outbound/models"
	org "euromoby.com/smsgw/internal/outbound/repos/group"
	orm "euromoby.com/smsgw/internal/outbound/repos/message"
	"euromoby.com/smsgw/internal/providers/connectors"

	coremodels "euromoby.com/core/models"
	corerepos "euromoby.com/core/repos"
)

const (
	defaultSender = "SMSGW"
	maxAttempts   = 5
)

type MessageWorker struct {
	app            *config.App
	connectorsRepo *connectors.ConnectorRepository
	b              billing.Billing
}

func NewMessageWorker(
	app *config.App,
	connectorsRepo *connectors.ConnectorRepository,
	b billing.Billing,
) *MessageWorker {
	return &MessageWorker{
		app:            app,
		connectorsRepo: connectorsRepo,
		b:              b,
	}
}

func (w *MessageWorker) Run() (bool, error) {
	tx, err := corerepos.Begin(w.app.Config.DBPool)
	if err != nil {
		return false, err
	}

	defer corerepos.Rollback(tx)

	message, err := orm.NewRepo(tx).FindOneForProcessing()
	if err != nil {
		if errors.Is(err, coremodels.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	logger.Infow("worker found new message for processing", "worker", w.Name(), "sms", message)

	if err := w.processMessage(tx, message); err != nil {
		return false, err
	}

	if err := corerepos.Commit(tx); err != nil {
		return false, err
	}

	return true, nil
}

func (w *MessageWorker) processMessage(tx db.Conn, message *models.Message) error {
	outboundMessageRepo := orm.NewRepo(tx)

	messageGroup, err := org.NewRepo(tx).FindByID(message.MessageGroupID)
	if err != nil {
		if errors.Is(err, coremodels.ErrNotFound) {
			logger.Errorw("message group is missing", "worker", w.Name(), "sms", message)
			message.Status = models.MessageStatusFailed

			return outboundMessageRepo.Update(message)
		}

		return err
	}

	w.sendToProvider(messageGroup, message)

	if err := outboundMessageRepo.Update(message); err != nil {
		return err
	}

	return w.sendNotification(tx, messageGroup, message)
}

func (w *MessageWorker) sendNotification(tx db.Conn, messageGroup *models.MessageGroup, message *models.Message) error {
	if message.Status == models.MessageStatusNew {
		return nil
	}

	if messageGroup.NotificationURL == nil {
		// TODO: check default callback
		return nil
	}

	n := nm.MakeOutboundDeliveryNotification(message)

	return notifications.NewRepo(tx).Save(n)
}

func (w *MessageWorker) sendToProvider(messageGroup *models.MessageGroup, message *models.Message) {
	connectorMessage := w.makeConnectorMessage(messageGroup, message)
	connector := w.connectorsRepo.FindConnector(connectorMessage)

	providerID := connector.Name()
	message.ProviderID = &providerID

	if err := w.b.ChargeOutboundMessage(message); err != nil {
		w.handleFailure(message, nil, err)

		return
	}

	resp, err := connector.SendMessage(connectorMessage)

	if err != nil {
		w.handleFailure(message, resp, err)
	} else {
		w.handleSuccess(message, resp)
	}
}

func (w *MessageWorker) handleFailure(message *models.Message, resp *connectors.SendMessageResponse, err error) {
	switch {
	case errors.Is(err, models.ErrSendFailed) || errors.Is(err, coremodels.ErrInvalidJSON):
		message.ProviderResponse = resp.Body
	case errors.Is(err, billing.ErrInsufficientFunds):
		message.Status = models.MessageStatusFailed
		errorText := err.Error()
		message.ProviderResponse = &errorText

		return
	case errors.Is(err, connectors.ErrDeadLetter):
		message.Status = models.MessageStatusFailed

		return
	default:
		errorText := err.Error()
		message.ProviderResponse = &errorText
	}

	w.tryReschedule(message)
}

func (w *MessageWorker) tryReschedule(message *models.Message) {
	if message.AttemptCounter >= w.MaxAttempts() {
		message.Status = models.MessageStatusFailed

		logger.Errorw("sending failed", "worker", w.Name(), "sms", message)

		return
	}

	message.NextAttemptAt = coremodels.CalculateNextAttemptTime(message.AttemptCounter)
	message.AttemptCounter++
	logger.Infow("sending attempt failed, will try again later",
		"worker", w.Name(),
		"sms", message,
		"next_attempt_at", message.NextAttemptAt,
		"attempt_counter", message.AttemptCounter,
	)
}

func (w *MessageWorker) handleSuccess(message *models.Message, resp *connectors.SendMessageResponse) {
	message.Status = models.MessageStatusSent
	message.ProviderResponse = resp.Body
	message.ProviderMessageID = resp.MessageID

	logger.Infow("successfully sent", "worker", w.Name(), "sms", message)
}

func (w *MessageWorker) makeConnectorMessage(messageGroup *models.MessageGroup, message *models.Message) *connectors.SendMessageRequest {
	var messageSender string
	if messageGroup.Sender != nil {
		messageSender = *messageGroup.Sender
	} else {
		messageSender = defaultSender
	}

	return &connectors.SendMessageRequest{
		MSISDN:              message.MSISDN,
		Sender:              messageSender,
		Body:                messageGroup.Body,
		ClientTransactionID: message.ID,
	}
}

func (w *MessageWorker) Name() string {
	return "OutboundMessageWorker"
}

func (w *MessageWorker) MaxAttempts() int {
	return maxAttempts
}

func (w *MessageWorker) SleepTime() time.Duration {
	return w.app.WorkerSleep
}
