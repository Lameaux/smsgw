package message

type Request struct {
	Sender              string `json:"sender"`
	MSISDN              string `json:"msisdn"`
	Body                string `json:"body"`
	ClientTransactionID string `json:"client_transaction_id"`
}
