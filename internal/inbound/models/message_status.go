package models

type MessageStatus string

const (
	MessageStatusNew       MessageStatus = "n"
	MessageStatusFailed    MessageStatus = "f"
	MessageStatusSent      MessageStatus = "s"
	MessageStatusDelivered MessageStatus = "d"
)
