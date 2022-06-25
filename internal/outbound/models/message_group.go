package models

import (
	"time"

	coremodels "euromoby.com/core/models"
)

type MessageGroup struct {
	ID                  string    `json:"id"`
	MerchantID          string    `json:"-"`
	Sender              *string   `json:"sender,omitempty"`
	Body                string    `json:"body"`
	NotificationURL     *string   `json:"notification_url,omitempty"`
	ClientTransactionID *string   `json:"client_transaction_id,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func NewMessageGroup(merchantID string, body string) *MessageGroup {
	generatedTransactionID := coremodels.NewUUID()
	now := coremodels.TimeNow()

	return &MessageGroup{
		MerchantID:          merchantID,
		Body:                body,
		ClientTransactionID: &generatedTransactionID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}
