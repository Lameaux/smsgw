package billing

import "euromoby.com/smsgw/internal/models"

type Billing interface {
	CheckBalance(merchantID string) error
	ChargeOutboundMessage(message *models.OutboundMessage) error
}
