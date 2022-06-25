package models

import (
	"time"

	coremodels "github.com/Lameaux/core/models"
)

type Message struct {
	ID                string            `json:"id"`
	MerchantID        string            `json:"-"`
	Shortcode         string            `json:"shortcode"`
	Status            MessageStatus     `json:"status"`
	MSISDN            coremodels.MSISDN `json:"msisdn"`
	Body              string            `json:"body"`
	ProviderID        string            `json:"-"`
	ProviderMessageID string            `json:"-"`
	NextAttemptAt     time.Time         `json:"-"`
	AttemptCounter    int               `json:"-"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}
