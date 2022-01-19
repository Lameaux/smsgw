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
	name          = "OutboundMessageWorker"
)

type OutboundMessageWorker struct {
	app            *config.AppConfig
	connectorsRepo *connectors.ConnectorRepository
}

func NewOutboundMessageWorker(app *config.AppConfig, connectorsRepo *connectors.ConnectorRepository) *OutboundMessageWorker {
	worker := OutboundMessageWorker{
		app:            app,
		connectorsRepo: connectorsRepo,
	}

	return &worker
}

func (w *OutboundMessageWorker) Run() (bool, error) {
	ctx, done := repos.DBTxContext()
	defer done()

	tx, err := w.app.DBPool.Begin(ctx)
	if err != nil {
		return false, err
	}

	outboundMessageRepo := repos.NewOutboundMessageRepo(tx)

	message, err := outboundMessageRepo.FindOneForProcessing()

	if err != nil {
		tx.Rollback(ctx)
		return false, err
	}

	if message == nil {
		tx.Rollback(ctx)
		return false, nil
	}

	logger.Infow("worker found new message for processing", "worker", w.Name(), "message", message)

	messageOrderRepo := repos.NewMessageOrderRepo(tx)

	messageOrder, err := messageOrderRepo.FindByMerchantAndID(message.MerchantID, message.MessageOrderID)
	if err != nil {
		tx.Rollback(ctx)
		return false, err
	}

	if messageOrder == nil {
		tx.Rollback(ctx)
		return false, nil
	}

	connectorMessage := w.makeConnectorMessage(messageOrder, message)
	connector := w.connectorsRepo.FindConnector(connectorMessage)
	resp, err := connector.SendMessage(connectorMessage)
	name := connector.Name()
	message.ProviderID = &name
	message.ProviderResponse = resp.Body

	if err != nil {
		message.Status = models.OutboundMessageStatusFailed
	} else {
		message.Status = models.OutboundMessageStatusSent
		message.ProviderMessageID = resp.MessageID
	}

	message.UpdatedAt = utils.Now()

	if err = outboundMessageRepo.UpdateStatus(message); err != nil {
		tx.Rollback(ctx)
		return false, err
	}

	tx.Commit(ctx)

	return true, nil
}

func (w *OutboundMessageWorker) makeConnectorMessage(messageOrder *models.MessageOrder, message *models.OutboundMessage) *connectors.MessageRequest {
	var messageSender string
	if messageOrder.Sender != nil {
		messageSender = *messageOrder.Sender
	} else {
		messageSender = defaultSender
	}

	return &connectors.MessageRequest{
		MSISDN: message.MSISDN,
		Sender: messageSender,
		Body:   messageOrder.Body,
	}
}

func (w *OutboundMessageWorker) Name() string {
	return name
}

func (w *OutboundMessageWorker) SleepTime() time.Duration {
	return w.app.WorkerSleep
}
