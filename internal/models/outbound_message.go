package models

import (
	"time"
)

type OutboundMessageStatus string

const (
	OutboundMessageStatusNew       OutboundMessageStatus = "n"
	OutboundMessageStatusFailed    OutboundMessageStatus = "f"
	OutboundMessageStatusSent      OutboundMessageStatus = "s"
	OutboundMessageStatusDelivered OutboundMessageStatus = "d"
)

type OutboundMessage struct {
	ID                string                `json:"id"`
	MerchantID        string                `json:"-"`
	MessageOrderID    string                `json:"message_order_id"`
	Status            OutboundMessageStatus `json:"status"`
	MSISDN            MSISDN                `json:"msisdn"`
	ProviderID        *string               `json:"-"`
	ProviderMessageID *string               `json:"-"`
	ProviderResponse  *string               `json:"-"`
	NextAttemptAt     time.Time             `json:"-"`
	AttemptCounter    int                   `json:"-"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

func NewOutboundMessage(merchantID string, messageOrderID string, msisdn MSISDN) *OutboundMessage {
	now := TimeNow()
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
