package billing

import om "github.com/Lameaux/smsgw/internal/outbound/models"

type Billing interface {
	CheckBalance(merchantID string) error
	Charge(message *om.Message) error
}
