package models

import (
	im "euromoby.com/smsgw/internal/inbound/models"
	om "euromoby.com/smsgw/internal/outbound/models"
	"time"

	coremodels "euromoby.com/core/models"
)

type DeliveryNotification struct {
	ID             string
	MessageType    DeliveryNotificationType
	MessageID      string
	Status         DeliveryNotificationStatus
	LastResponse   *string
	NextAttemptAt  time.Time
	AttemptCounter int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func MakeInboundDeliveryNotification(message *im.Message) *DeliveryNotification {
	now := coremodels.TimeNow()

	return &DeliveryNotification{
		MessageType:    DeliveryNotificationInbound,
		MessageID:      message.ID,
		Status:         DeliveryNotificationStatusNew,
		NextAttemptAt:  now,
		AttemptCounter: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func MakeOutboundDeliveryNotification(message *om.Message) *DeliveryNotification {
	now := coremodels.TimeNow()

	return &DeliveryNotification{
		MessageType:    DeliveryNotificationOutbound,
		MessageID:      message.ID,
		Status:         DeliveryNotificationStatusNew,
		NextAttemptAt:  now,
		AttemptCounter: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
