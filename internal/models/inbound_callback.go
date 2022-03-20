package models

import (
	"time"
)

type InboundCallback struct {
	ID        string    `json:"id"`
	Shortcode string    `json:"-"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewSimpleInboundCallback(shortcode string, url string) *InboundCallback {
	now := TimeNow()

	return &InboundCallback{
		Shortcode: shortcode,
		URL:       url,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
