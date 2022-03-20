package models

import (
	"time"
)

type InboundCallback struct {
	ID              string    `json:"id"`
	Shortcode       string    `json:"shortcode"`
	NotificationURL string    `json:"notification_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewSimpleInboundCallback(shortcode string, notificationURL string) *InboundCallback {
	now := TimeNow()

	return &InboundCallback{
		Shortcode:       shortcode,
		NotificationURL: notificationURL,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
