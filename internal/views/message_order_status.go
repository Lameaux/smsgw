package views

import (
	"fmt"

	"euromoby.com/smsgw/internal/models"
)

type MessageOrderDetail struct {
	models.MessageOrder

	Messages []*OutboundMessageDetail `json:"messages"`
	HREF     string                   `json:"href"`
}

func NewMessageOrderDetail(order *models.MessageOrder, messages []*models.OutboundMessage) *MessageOrderDetail {
	messageDetails := make([]*OutboundMessageDetail, 0, len(messages))

	for _, message := range messages {
		messageDetails = append(messageDetails, NewOutboundMessageDetail(message, nil))
	}

	return &MessageOrderDetail{
		MessageOrder: *order,
		Messages:     messageDetails,
		HREF:         fmt.Sprintf("/messages/status/%s", order.ID),
	}
}
