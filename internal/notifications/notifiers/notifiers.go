package notifiers

import "github.com/Lameaux/smsgw/internal/notifications/notifiers/models"

type Notifier interface {
	Send(url string, message interface{}) (*models.Response, error)
}
