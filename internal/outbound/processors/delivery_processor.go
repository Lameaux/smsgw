package processors

import (
	"github.com/Lameaux/core/db"
	nm "github.com/Lameaux/smsgw/internal/notifications/models"
	"github.com/Lameaux/smsgw/internal/notifications/notifiers"
	om "github.com/Lameaux/smsgw/internal/outbound/models"
	org "github.com/Lameaux/smsgw/internal/outbound/repos/group"
	orm "github.com/Lameaux/smsgw/internal/outbound/repos/message"
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
