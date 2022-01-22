package workers

import (
	"time"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/providers/connectors"
	"euromoby.com/smsgw/internal/repos"
	"euromoby.com/smsgw/internal/utils"
)

const (
	defaultSender = "SMSGW"
)

type OutboundMessageWorker struct {
	app            *config.AppConfig
	connectorsRepo *connectors.ConnectorRepository
}

func NewOutboundMessageWorker(app *config.AppConfig, connectorsRepo *connectors.ConnectorRepository) *OutboundMessageWorker {
	return &OutboundMessageWorker{
		app:            app,
		connectorsRepo: connectorsRepo,
	}
}

func (w *OutboundMessageWorker) Run() (bool, error) {
	ctx, done := repos.DBTxContext()
	defer done()

	tx, err := w.app.DBPool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	outboundMessageRepo := repos.NewOutboundMessageRepo(tx)

	message, err := outboundMessageRepo.FindOneForProcessing()
	if err != nil {
		return false, err
	}

	if message == nil {
		return false, nil
	}

	logger.Infow("worker found new message for processing", "worker", w.Name(), "sms", message)

	messageOrderRepo := repos.NewMessageOrderRepo(tx)

	messageOrder, err := messageOrderRepo.FindByID(message.MessageOrderID)
	if err != nil {
		return false, err
	}

	if messageOrder == nil {
		logger.Errorw("message order is missing", "worker", w.Name(), "sms", message)
		message.Status = models.OutboundMessageStatusFailed
	} else {
		w.sendToProvider(messageOrder, message)
	}

	if err = outboundMessageRepo.Update(message); err != nil {
		return false, err
	}

	if err = tx.Commit(ctx); err != nil {
		return false, nil
	}

	return true, nil
}

func (w *OutboundMessageWorker) sendToProvider(messageOrder *models.MessageOrder, message *models.OutboundMessage) {
	connectorMessage := w.makeConnectorMessage(messageOrder, message)
	connector := w.connectorsRepo.FindConnector(connectorMessage)
	resp, err := connector.SendMessage(connectorMessage)

	name := connector.Name()
	message.ProviderID = &name

	if err != nil {
		w.handleFailure(message, resp, err)
	} else {
		w.handleSuccess(message, resp)
	}
}

func (w *OutboundMessageWorker) handleFailure(message *models.OutboundMessage, resp *connectors.SendMessageResponse, err error) {
	switch err {
	case models.ErrSendFailed, models.ErrInvalidJSON:
		message.ProviderResponse = resp.Body
	default:
		errorText := err.Error()
		message.ProviderResponse = &errorText
	}

	w.tryReschedule(message)
}

func (w *OutboundMessageWorker) tryReschedule(message *models.OutboundMessage) {
	if message.AttemptCounter >= w.MaxAttempts() {
		message.Status = models.OutboundMessageStatusFailed
		logger.Errorw("sending failed", "worker", w.Name(), "sms", message)
		return
	}

	message.NextAttemptAt = utils.CalculateNextAttemptTime(message.AttemptCounter)
	message.AttemptCounter++
	logger.Infow("sending attempt failed, will try again later",
		"worker", w.Name(),
		"sms", message,
		"next_attempt_at", message.NextAttemptAt,
		"attempt_counter", message.AttemptCounter,
	)
}

func (w *OutboundMessageWorker) handleSuccess(message *models.OutboundMessage, resp *connectors.SendMessageResponse) {
	message.Status = models.OutboundMessageStatusSent
	message.ProviderResponse = resp.Body
	message.ProviderMessageID = resp.MessageID

	logger.Infow("successfully sent", "worker", w.Name(), "sms", message)
}

func (w *OutboundMessageWorker) makeConnectorMessage(messageOrder *models.MessageOrder, message *models.OutboundMessage) *connectors.SendMessageRequest {
	var messageSender string
	if messageOrder.Sender != nil {
		messageSender = *messageOrder.Sender
	} else {
		messageSender = defaultSender
	}

	return &connectors.SendMessageRequest{
		MSISDN:              message.MSISDN,
		Sender:              messageSender,
		Body:                messageOrder.Body,
		ClientTransactionID: message.ID,
	}
}

func (w *OutboundMessageWorker) Name() string {
	return "OutboundMessageWorker"
}

func (w *OutboundMessageWorker) MaxAttempts() int {
	return 5
}

func (w *OutboundMessageWorker) SleepTime() time.Duration {
	return w.app.WorkerSleep
}
