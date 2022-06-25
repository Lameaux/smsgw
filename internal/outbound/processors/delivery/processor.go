package delivery

import (
	"github.com/Lameaux/core/db"
	nm "github.com/Lameaux/smsgw/internal/notifications/models"
	"github.com/Lameaux/smsgw/internal/notifications/notifiers"
	notifiersmodels "github.com/Lameaux/smsgw/internal/notifications/notifiers/models"
	om "github.com/Lameaux/smsgw/internal/outbound/models"
	org "github.com/Lameaux/smsgw/internal/outbound/repos/group"
	orm "github.com/Lameaux/smsgw/internal/outbound/repos/message"
)

type Processor struct {
	notifier notifiers.Notifier
}

func NewProcessor(notifier notifiers.Notifier) *Processor {
	return &Processor{notifier}
}

func (p *Processor) SendNotification(db db.Conn, notification *nm.DeliveryNotification) (*notifiersmodels.Response, error) {
	message, messageGroup, err := findMessageWithGroup(db, notification.MessageID)
	if err != nil {
		return nil, err
	}

	if messageGroup.NotificationURL == nil {
		// TODO: check default callback
		return nil, nm.ErrMissingNotificationURL
	}

	return p.notifier.Send(*messageGroup.NotificationURL, message)
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
