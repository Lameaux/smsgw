package models

import (
	"time"

	"euromoby.com/smsgw/internal/utils"
)

type (
	MessageType                string
	DeliveryNotificationStatus string
)

const (
	MessageTypeOutbound MessageType = "o"
	MessageTypeInbound  MessageType = "i"
)

const (
	DeliveryNotificationStatusNew    DeliveryNotificationStatus = "n"
	DeliveryNotificationStatusFailed DeliveryNotificationStatus = "f"
	DeliveryNotificationStatusSent   DeliveryNotificationStatus = "s"
)

type DeliveryNotification struct {
	ID             string
	MessageType    MessageType
	MessageID      string
	Status         DeliveryNotificationStatus
	LastResponse   *string
	NextAttemptAt  time.Time
	AttemptCounter int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func MakeInboundDeliveryNotification(message *InboundMessage) *DeliveryNotification {
	now := utils.Now()
	return &DeliveryNotification{
		MessageType:    MessageTypeInbound,
		MessageID:      message.ID,
		Status:         DeliveryNotificationStatusNew,
		NextAttemptAt:  now,
		AttemptCounter: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
