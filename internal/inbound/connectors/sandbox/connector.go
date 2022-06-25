package sandbox

import (
	"github.com/Lameaux/smsgw/internal/config"
	"github.com/Lameaux/smsgw/internal/inbound/connectors"
	is "github.com/Lameaux/smsgw/internal/inbound/connectors/inputs/status"
	os "github.com/Lameaux/smsgw/internal/inbound/connectors/outputs/status"
)

const (
	apiBaseURL = "http://0.0.0.0:8081/sandbox"
)

type Connector struct {
	app *config.App
}

func NewConnector(app *config.App) *Connector {
	return &Connector{
		app: app,
	}
}

func (c *Connector) Name() string {
	return "sandbox"
}

func (c *Connector) Accept(status *is.Request) bool {
	return status.Provider == c.Name()
}

func (c *Connector) SendStatus(status *is.Request) (*os.Response, error) {
	return nil, connectors.ErrSendFailed
}
