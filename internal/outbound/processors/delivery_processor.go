package processors

import (
	"euromoby.com/core/db"
	nm "euromoby.com/smsgw/internal/notifications/models"
	"euromoby.com/smsgw/internal/notifications/notifiers"
	om "euromoby.com/smsgw/internal/outbound/models"
	org "euromoby.com/smsgw/internal/outbound/repos/group"
	orm "euromoby.com/smsgw/internal/outbound/repos/message"
)

type DeliveryProcessor struct {
	notifier *notifiers.DefaultNotifier
}

func NewDeliveryProcessor(notifier *notifiers.DefaultNotifier) *DeliveryProcessor {
	return &DeliveryProcessor{notifier}
}

func (p *DeliveryProcessor) SendNotification(db db.Conn, notification *nm.DeliveryNotification) (*notifiers.NotifierResponse, error) {
	message, messageGroup, err := findMessageWithGroup(db, notification.MessageID)
	if err != nil {
		return nil, err
	}

	if messageGroup.NotificationURL == nil {
		// TODO: check default callback
		return nil, nm.ErrMissingNotificationURL
	}

	return p.notifier.SendNotification(*messageGroup.NotificationURL, message)
}

func findMessageWithGroup(db db.Conn, id string) (*om.Message, *om.MessageGroup, error) {
	message, err := orm.NewRepo(db).FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	messageGroup, err := org.NewRepo(db).FindByID(message.MessageGroupID)
	if err != nil {
		return nil, nil, err
	}

	return message, messageGroup, nil
}
