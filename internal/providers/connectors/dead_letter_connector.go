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

func (c *DeadLetterConnector) Accept(message *SendMessageRequest) bool {
	return true
}

func (c *DeadLetterConnector) SendMessage(message *SendMessageRequest) (*SendMessageResponse, error) {
	return nil, models.ErrDeadLetter
}

func (c *DeadLetterConnector) SendStatus(status *SendStatusRequest) (*SendStatusResponse, error) {
	return nil, models.ErrDeadLetter
}
