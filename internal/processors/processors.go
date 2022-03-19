package processors

import (
	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/notifiers"
)

type DeliveryNotificationProcessor interface {
	SendNotification(db db.Conn, notification *models.DeliveryNotification) (*notifiers.NotifierResponse, error)
}
