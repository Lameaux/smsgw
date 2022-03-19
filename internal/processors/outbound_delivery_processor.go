package processors

import (
	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/notifiers"
	"euromoby.com/smsgw/internal/repos"
)

type OutboundDeliveryProcessor struct {
	notifier *notifiers.DefaultNotifier
}

func NewOutboundDeliveryProcessor(notifier *notifiers.DefaultNotifier) *OutboundDeliveryProcessor {
	return &OutboundDeliveryProcessor{notifier}
}

func (p *OutboundDeliveryProcessor) SendNotification(db db.Conn, notification *models.DeliveryNotification) (*notifiers.NotifierResponse, error) {
	message, messageOrder, err := p.findMessageWithOrder(db, notification.MessageID)
	if err != nil {
		return nil, err
	}

	if messageOrder.NotificationURL == nil {
		return nil, models.ErrMissingNotificationURL
	}

	return p.notifier.SendNotification(*messageOrder.NotificationURL, message)
}

func (p *OutboundDeliveryProcessor) findMessageWithOrder(db db.Conn, id string) (*models.OutboundMessage, *models.MessageOrder, error) {
	outboundMessageRepo := repos.NewOutboundMessageRepo(db)

	message, err := outboundMessageRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	messageOrderRepo := repos.NewMessageOrderRepo(db)

	messageOrder, err := messageOrderRepo.FindByID(message.MessageOrderID)
	if err != nil {
		return nil, nil, err
	}

	return message, messageOrder, nil
}
