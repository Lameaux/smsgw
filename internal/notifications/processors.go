package notifications

import (
	"github.com/Lameaux/core/db"
	"github.com/Lameaux/smsgw/internal/notifications/models"
	nm "github.com/Lameaux/smsgw/internal/notifications/notifiers/models"
)

type Processor interface {
	SendNotification(db db.Conn, notification *models.DeliveryNotification) (*nm.Response, error)
}
