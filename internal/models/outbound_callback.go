package models

import (
	"time"
)

type OutboundCallback struct {
	ID              string    `json:"id"`
	MerchantID      string    `json:"-"`
	NotificationURL string    `json:"notification_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewSimpleOutboundCallback(merchantID string, notificationURL string) *OutboundCallback {
	now := TimeNow()

	return &OutboundCallback{
		MerchantID:      merchantID,
		NotificationURL: notificationURL,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
