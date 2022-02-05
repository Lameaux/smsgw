package inputs

import "euromoby.com/smsgw/internal/models"

type SendMessageParams struct {
	MerchantID          string
	To                  []string        `json:"to"`
	Recipients          []models.MSISDN `json:"-"`
	Sender              *string         `json:"sender,omitempty"`
	Body                string          `json:"body"`
	NotificationURL     *string         `json:"notification_url,omitempty"`
	ClientTransactionID *string         `json:"client_transaction_id,omitempty"`
}
