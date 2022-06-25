package connectors

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
	return nil, ErrDeadLetter
}

func (c *DeadLetterConnector) SendStatus(status *SendStatusRequest) (*SendStatusResponse, error) {
	return nil, ErrDeadLetter
}
