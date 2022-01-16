package models

import (
	"time"

	"euromoby.com/smsgw/internal/utils"
)

type OutboundMessageStatus string

const (
	OutboundMessageStatusNew       OutboundMessageStatus = "new"
	OutboundMessageStatusFailed    OutboundMessageStatus = "failed"
	OutboundMessageStatusSent      OutboundMessageStatus = "sent"
	OutboundMessageStatusDelivered OutboundMessageStatus = "delivered"
)

type OutboundMessage struct {
	ID                string                `json:"id"`
	MerchantID        string                `json:"-"`
	MessageOrderID    string                `json:"message_order_id"`
	Status            OutboundMessageStatus `json:"status"`
	MSISDN            string                `json:"msisdn"`
	ProviderID        *string               `json:"-"`
	ProviderMessageID *string               `json:"-"`
	ProviderResponse  *string               `json:"-"`
	NextAttemptAt     time.Time             `json:"-"`
	AttemptCounter    int                   `json:"-"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

func NewOutboundMessage(merchantID string, messageOrderID string, msisdn string) *OutboundMessage {
	now := utils.Now()
	return &OutboundMessage{
		MerchantID:     merchantID,
		MessageOrderID: messageOrderID,
		Status:         OutboundMessageStatusNew,
		MSISDN:         msisdn,
		NextAttemptAt:  now,
		AttemptCounter: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
