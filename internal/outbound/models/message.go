package models

import (
	"time"

	coremodels "euromoby.com/core/models"
)

type Message struct {
	ID                string            `json:"id"`
	MerchantID        string            `json:"-"`
	MessageGroupID    string            `json:"message_group_id"`
	Status            MessageStatus     `json:"status"`
	MSISDN            coremodels.MSISDN `json:"msisdn"`
	ProviderID        *string           `json:"-"`
	ProviderMessageID *string           `json:"-"`
	ProviderResponse  *string           `json:"-"`
	NextAttemptAt     time.Time         `json:"-"`
	AttemptCounter    int               `json:"-"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

func NewMessage(merchantID string, messageGroupID string, msisdn coremodels.MSISDN) *Message {
	now := coremodels.TimeNow()

	return &Message{
		MerchantID:     merchantID,
		MessageGroupID: messageGroupID,
		Status:         MessageStatusNew,
		MSISDN:         msisdn,
		NextAttemptAt:  now,
		AttemptCounter: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
