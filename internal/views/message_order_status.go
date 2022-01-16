package views

import "euromoby.com/smsgw/internal/models"

type MessageOrderStatus struct {
	models.MessageOrder

	Messages []*models.OutboundMessage `json:"messages"`
}

func NewMessageOrderStatus(order *models.MessageOrder, messages []*models.OutboundMessage) *MessageOrderStatus {
	return &MessageOrderStatus{
		MessageOrder: *order,
		Messages:     messages,
	}
}
