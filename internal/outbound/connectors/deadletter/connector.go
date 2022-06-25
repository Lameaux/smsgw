package deadletter

import (
	im "github.com/Lameaux/smsgw/internal/outbound/connectors/inputs/message"
	"github.com/Lameaux/smsgw/internal/outbound/connectors/models"
	om "github.com/Lameaux/smsgw/internal/outbound/connectors/outputs/message"
)

type Connector struct{}

func NewConnector() *Connector {
	return &Connector{}
}

func (c *Connector) Name() string {
	return "deadletter"
}

func (c *Connector) Accept(message *im.Request) bool {
	return true
}

func (c *Connector) SendMessage(message *im.Request) (*om.Response, error) {
	return nil, models.ErrDeadLetter
}
