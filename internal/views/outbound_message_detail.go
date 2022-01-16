package views

import "euromoby.com/smsgw/internal/models"

type OutboundMessageDetail struct {
	models.OutboundMessage
	MessageOrder models.MessageOrder `json:"message_order"`
}

func NewOutboundMessageDetail(message *models.OutboundMessage, order *models.MessageOrder) *OutboundMessageDetail {
	return &OutboundMessageDetail{
		OutboundMessage: *message,
		MessageOrder:    *order,
	}
}
