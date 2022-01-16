package sandbox

type SandboxInboundMessage struct {
	MessageID string `json:"message_id"`
	Shortcode string `json:"shortcode"`
	MSISDN    string `json:"msisdn"`
	Body      string `json:"body"`
}
