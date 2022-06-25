package notifications

import (
	"github.com/Lameaux/core/db"
	"github.com/Lameaux/smsgw/internal/notifications/models"
	"github.com/Lameaux/smsgw/internal/notifications/notifiers"
)

type Processor interface {
	SendNotification(db db.Conn, notification *models.DeliveryNotification) (*notifiers.NotifierResponse, error)
}
