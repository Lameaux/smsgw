package models

type DeliveryNotificationStatus string

const (
	DeliveryNotificationStatusNew    DeliveryNotificationStatus = "n"
	DeliveryNotificationStatusFailed DeliveryNotificationStatus = "f"
	DeliveryNotificationStatusSent   DeliveryNotificationStatus = "s"
)
