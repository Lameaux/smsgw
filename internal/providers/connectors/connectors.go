package connectors

import "euromoby.com/smsgw/internal/logger"

type SendMessageRequest struct {
	MSISDN              string
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
	MSISDN              string
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

func NewConnectorRepository() *ConnectorRepository {
	connectors := make(map[string]Connector)

	sandboxConnector := NewSandboxConnector()
	connectors[sandboxConnector.Name()] = sandboxConnector

	return &ConnectorRepository{connectors}
}

func (r *ConnectorRepository) FindConnectorByName(name string) (Connector, bool) {
	connector, found := r.connectors[name]
	return connector, found
}

func (r *ConnectorRepository) FindConnector(message *SendMessageRequest) Connector {
	for _, connector := range r.connectors {
		if connector.Accept(message) {
			return connector
		}
	}

	logger.Infow("no connector found for the message", "sms", message)
	return &DeadLetterConnector{}
}
