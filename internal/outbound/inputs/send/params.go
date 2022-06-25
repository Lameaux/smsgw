package send

import coremodels "github.com/Lameaux/core/models"

type Params struct {
	MerchantID          string
	To                  []string            `json:"to"`
	Recipients          []coremodels.MSISDN `json:"-"`
	Sender              *string             `json:"sender,omitempty"`
	Body                string              `json:"body"`
	NotificationURL     *string             `json:"notification_url,omitempty"`
	ClientTransactionID *string             `json:"client_transaction_id,omitempty"`
}
