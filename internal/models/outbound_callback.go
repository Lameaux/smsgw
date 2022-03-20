package models

import (
	"time"
)

type OutboundCallback struct {
	ID         string    `json:"id"`
	MerchantID string    `json:"-"`
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewSimpleOutboundCallback(merchantID string, url string) *OutboundCallback {
	now := TimeNow()

	return &OutboundCallback{
		MerchantID: merchantID,
		URL:        url,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
