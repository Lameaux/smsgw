package models

type (
	DeliveryNotificationType string
)

const (
	DeliveryNotificationOutbound DeliveryNotificationType = "o"
	DeliveryNotificationInbound  DeliveryNotificationType = "i"
)
