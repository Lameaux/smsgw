package connectors

import (
	is "github.com/Lameaux/smsgw/internal/inbound/connectors/inputs/status"
	os "github.com/Lameaux/smsgw/internal/inbound/connectors/outputs/status"
)

type Connector interface {
	Name() string
	Accept(message *is.Request) bool
	SendStatus(message *is.Request) (*os.Response, error)
}
