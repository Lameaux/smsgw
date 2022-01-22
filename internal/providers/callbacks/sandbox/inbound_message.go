package sandbox

type InboundMessage struct {
	MessageID string `json:"message_id"`
	Shortcode string `json:"shortcode"`
	MSISDN    string `json:"msisdn"`
	Body      string `json:"body"`
}
