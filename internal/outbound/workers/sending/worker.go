package sending

import (
	"errors"
	"github.com/Lameaux/smsgw/internal/billing"
	"github.com/Lameaux/smsgw/internal/notifications"
	"time"

	"github.com/Lameaux/core/db"
	"github.com/Lameaux/core/logger"
	"github.com/Lameaux/smsgw/internal/config"
	nm "github.com/Lameaux/smsgw/internal/notifications/models"
	"github.com/Lameaux/smsgw/internal/outbound/connectors"
	cm "github.com/Lameaux/smsgw/internal/outbound/connectors/models"
	"github.com/Lameaux/smsgw/internal/outbound/models"
	org "github.com/Lameaux/smsgw/internal/outbound/repos/group"
	orm "github.com/Lameaux/smsgw/internal/outbound/repos/message"

	cim "github.com/Lameaux/smsgw/internal/outbound/connectors/inputs/message"
	com "github.com/Lameaux/smsgw/internal/outbound/connectors/outputs/message"

	coremodels "github.com/Lameaux/core/models"
	corerepos "github.com/Lameaux/core/repos"
)

const (
	defaultSender = "SMSGW"
	maxAttempts   = 5
)

type Worker struct {
	app            *config.App
	connectorsRepo *connectors.ConnectorRepository
	b              billing.Billing
}

func NewWorker(
	app *config.App,
	connectorsRepo *connectors.ConnectorRepository,
	b billing.Billing,
) *Worker {
	return &Worker{
		app:            app,
		connectorsRepo: connectorsRepo,
		b:              b,
	}
}

func (w *Worker) Run() (bool, error) {
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

func (w *Worker) processMessage(tx db.Conn, message *models.Message) error {
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

func (w *Worker) sendNotification(tx db.Conn, messageGroup *models.MessageGroup, message *models.Message) error {
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

func (w *Worker) sendToProvider(messageGroup *models.MessageGroup, message *models.Message) {
	connectorMessage := w.makeConnectorMessage(messageGroup, message)
	connector := w.connectorsRepo.FindConnector(connectorMessage)

	providerID := connector.Name()
	message.ProviderID = &providerID

	if err := w.b.Charge(message); err != nil {
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

func (w *Worker) handleFailure(message *models.Message, resp *com.Response, err error) {
	switch {
	case errors.Is(err, models.ErrSendFailed) || errors.Is(err, coremodels.ErrInvalidJSON):
		message.ProviderResponse = resp.Body
	case errors.Is(err, billing.ErrInsufficientFunds):
		message.Status = models.MessageStatusFailed
		errorText := err.Error()
		message.ProviderResponse = &errorText

		return
	case errors.Is(err, cm.ErrDeadLetter):
		message.Status = models.MessageStatusFailed

		return
	default:
		errorText := err.Error()
		message.ProviderResponse = &errorText
	}

	w.tryReschedule(message)
}

func (w *Worker) tryReschedule(message *models.Message) {
	if message.AttemptCounter >= MaxAttempts() {
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

func (w *Worker) handleSuccess(message *models.Message, resp *com.Response) {
	message.Status = models.MessageStatusSent
	message.ProviderResponse = resp.Body
	message.ProviderMessageID = resp.MessageID

	logger.Infow("successfully sent", "worker", w.Name(), "sms", message)
}

func (w *Worker) makeConnectorMessage(messageGroup *models.MessageGroup, message *models.Message) *cim.Request {
	var messageSender string
	if messageGroup.Sender != nil {
		messageSender = *messageGroup.Sender
	} else {
		messageSender = defaultSender
	}

	return &cim.Request{
		MSISDN:              message.MSISDN,
		Sender:              messageSender,
		Body:                messageGroup.Body,
		ClientTransactionID: message.ID,
	}
}

func (w *Worker) Name() string {
	return "OutboundMessageWorker"
}

func MaxAttempts() int {
	return maxAttempts
}

func (w *Worker) SleepTime() time.Duration {
	return w.app.WorkerSleep
}
