package models

import "time"

type InboundMessageStatus string

const (
	InboundMessageStatusNew       InboundMessageStatus = "n"
	InboundMessageStatusFailed    InboundMessageStatus = "f"
	InboundMessageStatusSent      InboundMessageStatus = "s"
	InboundMessageStatusDelivered InboundMessageStatus = "d"
)

type InboundMessage struct {
	ID                string               `json:"id"`
	Shortcode         string               `json:"shortcode"`
	Status            InboundMessageStatus `json:"status"`
	MSISDN            MSISDN               `json:"msisdn"`
	Body              string               `json:"body"`
	ProviderID        string               `json:"-"`
	ProviderMessageID string               `json:"-"`
	NextAttemptAt     time.Time            `json:"-"`
	AttemptCounter    int                  `json:"-"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}
