package views

type SendMessageParams struct {
	To                  []string `json:"to"`
	Sender              *string  `json:"sender,omitempty"`
	Body                string   `json:"body"`
	NotificationURL     *string  `json:"notification_url,omitempty"`
	ClientTransactionID *string  `json:"client_transaction_id,omitempty"`
}
