package models

import "time"

type InboundMessageStatus string

const (
	InboundMessageStatusNew InboundMessageStatus = "new"
	InboundMessageStatusAck InboundMessageStatus = "ack"
)

type InboundMessage struct {
	ID                string               `json:"id"`
	Shortcode         string               `json:"shortcode"`
	Status            InboundMessageStatus `json:"status"`
	MSISDN            string               `json:"msisdn"`
	Body              string               `json:"body"`
	ProviderID        string               `json:"-"`
	ProviderMessageID string               `json:"-"`
	NextAttemptAt     time.Time            `json:"-"`
	AttemptCounter    int                  `json:"-"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}
