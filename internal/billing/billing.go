package billing

import om "euromoby.com/smsgw/internal/outbound/models"

type Billing interface {
	CheckBalance(merchantID string) error
	ChargeOutboundMessage(message *om.Message) error
}
