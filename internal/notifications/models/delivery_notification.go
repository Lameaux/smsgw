package models

import (
	im "github.com/Lameaux/smsgw/internal/inbound/models"
	om "github.com/Lameaux/smsgw/internal/outbound/models"
	"time"

	coremodels "github.com/Lameaux/core/models"
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
