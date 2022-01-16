package connectors

import "euromoby.com/smsgw/internal/logger"

type MessageRequest struct {
	MSISDN string
	Sender string
	Body   string
}

type MessageResponse struct {
	MessageID *string
	Body      *string
}

type StatusRequest struct {
	MessageID string
	MSISDN    string
	Status    string
}

type StatusResponse struct {
	Body *string
}

type Connector interface {
	Name() string
	Accept(message *MessageRequest) bool
	SendMessage(message *MessageRequest) (*MessageResponse, error)
	SendStatus(message *StatusRequest) (*StatusResponse, error)
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

func (r *ConnectorRepository) FindConnector(message *MessageRequest) Connector {
	for _, connector := range r.connectors {
		if connector.Accept(message) {
			return connector
		}
	}

	logger.Infow("no connector found for the message", "message", message)
	return &DeadLetterConnector{}
}
