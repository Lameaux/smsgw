package models

import (
	"time"
)

type InboundCallback struct {
	ID         string    `json:"-"`
	MerchantID string    `json:"-"`
	Shortcode  string    `json:"-"`
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewSimpleInboundCallback(merchantID string, shortcode string, url string) *InboundCallback {
	now := TimeNow()

	return &InboundCallback{
		MerchantID: merchantID,
		Shortcode:  shortcode,
		URL:        url,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
