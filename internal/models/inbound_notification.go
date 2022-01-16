package models

import (
	"time"

	"euromoby.com/smsgw/internal/utils"
)

type (
	InboundNotificationType   string
	InboundNotificationStatus string
)

const (
	InboundNotificationStatusNew       InboundNotificationStatus = "new"
	InboundNotificationStatusFailed    InboundNotificationStatus = "failed"
	InboundNotificationStatusDelivered InboundNotificationStatus = "delivered"
)

type InboundNotification struct {
	ID               string
	MessageID        string
	Status           InboundNotificationStatus
	ProviderResponse *string
	NextAttemptAt    time.Time
	AttemptCounter   int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func MakeInboundNotification(message *InboundMessage) *InboundNotification {
	now := utils.Now()
	return &InboundNotification{
		MessageID:      message.ID,
		Status:         InboundNotificationStatusNew,
		NextAttemptAt:  now,
		AttemptCounter: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
