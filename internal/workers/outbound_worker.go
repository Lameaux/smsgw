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

type OutboundWorker struct {
	app            *config.AppConfig
	connectorsRepo *connectors.ConnectorRepository
	done           chan bool
}

func NewOutboundWorker(app *config.AppConfig, connectorsRepo *connectors.ConnectorRepository) (*OutboundWorker, error) {
	worker := OutboundWorker{
		app:            app,
		connectorsRepo: connectorsRepo,
		done:           make(chan bool, 1),
	}

	return &worker, nil
}

func (w *OutboundWorker) Run() {
	logger.Infow("OutboundWorker started")
	for {
		select {
		case <-w.done:
			logger.Infow("OutboundWorker stopped")
			return
		default:
			w.workAndSleep()
		}
	}
}

func (w *OutboundWorker) process() {

	for {
		found, err := w.processOneRecord()

		if err != nil {
			logger.Error(err)
			return
		}
		if !found {
			logger.Infow("OutboundWorker found nothing to process")
			return
		}
	}
}

func (w *OutboundWorker) processOneRecord() (bool, error) {
	ctx, done := repos.DBTxContext()
	defer done()

	tx, err := w.app.DBPool.Begin(ctx)
	if err != nil {
		return false, err
	}

	outboundMessageRepo := repos.NewOutboundMessageRepo(tx)

	message, err := outboundMessageRepo.FindOneForProcessing(models.OutboundMessageStatusNew)

	if err != nil {
		tx.Rollback(ctx)
		return false, err
	}

	if message == nil {
		tx.Rollback(ctx)
		return false, nil
	}

	logger.Infow("OutboundWorker found new message for processing", "message", message)

	messageOrderRepo := repos.NewMessageOrderRepo(tx)

	messageOrder, err := messageOrderRepo.FindByID(message.MerchantID, message.MessageOrderID)
	if err != nil {
		tx.Rollback(ctx)
		return false, err
	}

	if messageOrder == nil {
		tx.Rollback(ctx)
		return false, nil
	}

	connectorMessage := makeConnectorMessage(messageOrder, message)
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

func (w *OutboundWorker) workAndSleep() {
	logger.Infow("OutboundWorker working")
	w.process()
	logger.Infow("OutboundWorker sleeping")
	time.Sleep(w.app.WorkerSleep)
}

func (w *OutboundWorker) Terminate() {
	close(w.done)
}

func makeConnectorMessage(messageOrder *models.MessageOrder, message *models.OutboundMessage) *connectors.MessageRequest {
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
