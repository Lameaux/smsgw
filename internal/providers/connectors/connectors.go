package connectors

import (
	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/models"
)

type SendMessageRequest struct {
	MSISDN              models.MSISDN
	Sender              string
	Body                string
	ClientTransactionID string
}

type SendMessageResponse struct {
	MessageID *string
	Body      *string
}

type SendStatusRequest struct {
	MessageID           string
	MSISDN              models.MSISDN
	Status              string
	ClientTransactionID string
}

type SendStatusResponse struct {
	Body *string
}

type Connector interface {
	Name() string
	Accept(message *SendMessageRequest) bool
	SendMessage(message *SendMessageRequest) (*SendMessageResponse, error)
	SendStatus(message *SendStatusRequest) (*SendStatusResponse, error)
}

type ConnectorRepository struct {
	connectors map[string]Connector
}

func NewConnectorRepository(app *config.AppConfig) *ConnectorRepository {
	connectors := make(map[string]Connector)

	sandboxConnector := NewSandboxConnector(app)
	connectors[sandboxConnector.Name()] = sandboxConnector

	return &ConnectorRepository{connectors}
}

func (r *ConnectorRepository) FindConnectorByName(name string) (Connector, bool) { //nolint:ireturn
	connector, found := r.connectors[name]

	return connector, found
}

func (r *ConnectorRepository) FindConnector(message *SendMessageRequest) Connector { //nolint:ireturn
	for _, connector := range r.connectors {
		if connector.Accept(message) {
			return connector
		}
	}

	logger.Infow("no connector found for the message", "sms", message)

	return &DeadLetterConnector{}
}
