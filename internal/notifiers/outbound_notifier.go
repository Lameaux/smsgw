package notifiers

import "euromoby.com/smsgw/internal/models"

type OutboundNotifier struct{}

func NewOutboundNotifier() *OutboundNotifier {
	return &OutboundNotifier{}
}

func (*OutboundNotifier) SendNotification(messageOrder *models.MessageOrder, message *models.OutboundMessage) (*SendNotificationResponse, error) {
	body := "error"
	r := SendNotificationResponse{
		Body: &body,
	}
	return &r, models.ErrSendFailed
}
