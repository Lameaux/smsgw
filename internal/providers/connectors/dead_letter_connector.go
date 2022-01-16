package connectors

import (
	"euromoby.com/smsgw/internal/models"
)

type DeadLetterConnector struct {
	name string
}

func NewDeadLetterConnector() *DeadLetterConnector {
	return &DeadLetterConnector{
		name: "deadletter",
	}
}

func (c *DeadLetterConnector) Name() string {
	return c.name
}

func (c *DeadLetterConnector) Accept(message *MessageRequest) bool {
	return true
}

func (c *DeadLetterConnector) SendMessage(message *MessageRequest) (*MessageResponse, error) {
	return nil, models.ErrDeadLetter
}

func (c *DeadLetterConnector) SendStatus(status *StatusRequest) (*StatusResponse, error) {
	return nil, models.ErrDeadLetter
}
