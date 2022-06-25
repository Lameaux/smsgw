package notifications

import (
	"euromoby.com/core/db"
	"euromoby.com/smsgw/internal/notifications/models"
	"euromoby.com/smsgw/internal/notifications/notifiers"
)

type Processor interface {
	SendNotification(db db.Conn, notification *models.DeliveryNotification) (*notifiers.NotifierResponse, error)
}
