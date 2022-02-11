package models

import (
	"time"
)

type MessageOrder struct {
	ID                  string    `json:"id"`
	MerchantID          string    `json:"-"`
	Sender              *string   `json:"sender,omitempty"`
	Body                string    `json:"body"`
	NotificationURL     *string   `json:"notification_url,omitempty"`
	ClientTransactionID *string   `json:"client_transaction_id,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func NewSimpleMessageOrder(merchantID string, body string) *MessageOrder {
	generatedTransactionID := NewUUID()
	now := TimeNow()

	return &MessageOrder{
		MerchantID:          merchantID,
		Body:                body,
		ClientTransactionID: &generatedTransactionID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}
