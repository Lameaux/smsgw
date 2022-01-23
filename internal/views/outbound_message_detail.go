package views

import (
	"fmt"

	"euromoby.com/smsgw/internal/models"
)

type OutboundMessageDetail struct {
	models.OutboundMessage
	MessageOrder *MessageOrderDetail `json:"message_order,omitempty"`
	HREF         string              `json:"href"`
}

func NewOutboundMessageDetail(message *models.OutboundMessage, order *models.MessageOrder) *OutboundMessageDetail {
	var detail *MessageOrderDetail

	if order != nil {
		detail = NewMessageOrderDetail(order, nil)
	}

	return &OutboundMessageDetail{
		OutboundMessage: *message,
		MessageOrder:    detail,
		HREF:            fmt.Sprintf("/messages/outbound/%s", message.ID),
	}
}
