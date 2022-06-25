package connectors

import (
	"github.com/Lameaux/core/logger"
	"github.com/Lameaux/smsgw/internal/config"
	"github.com/Lameaux/smsgw/internal/outbound/connectors/deadletter"
	"github.com/Lameaux/smsgw/internal/outbound/connectors/sandbox"

	im "github.com/Lameaux/smsgw/internal/outbound/connectors/inputs/message"
	om "github.com/Lameaux/smsgw/internal/outbound/connectors/outputs/message"
)

type Connector interface {
	Name() string
	Accept(message *im.Request) bool
	SendMessage(message *im.Request) (*om.Response, error)
}

type ConnectorRepository struct {
	connectors map[string]Connector
}

func NewConnectorRepository(app *config.App) *ConnectorRepository {
	connectors := make(map[string]Connector)

	sandboxConnector := sandbox.NewConnector(app)
	connectors[sandboxConnector.Name()] = sandboxConnector

	return &ConnectorRepository{connectors}
}

func (r *ConnectorRepository) FindConnectorByName(name string) (Connector, bool) { //nolint:ireturn
	connector, found := r.connectors[name]

	return connector, found
}

func (r *ConnectorRepository) FindConnector(message *im.Request) Connector { //nolint:ireturn
	for _, connector := range r.connectors {
		if connector.Accept(message) {
			return connector
		}
	}

	logger.Infow("no connector found for the message", "sms", message)

	return deadletter.NewConnector()
}
