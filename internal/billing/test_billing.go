package billing

import (
	"euromoby.com/smsgw/internal/models"
)

type TestBilling struct{}

func NewTestBilling() *TestBilling {
	return &TestBilling{}
}

func (b *TestBilling) CheckBalance(merchantID string) error {
	return nil
}

func (b *TestBilling) ChargeOutboundMessage(message *models.OutboundMessage) error {
	return nil
}
