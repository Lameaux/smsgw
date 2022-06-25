package billing

import om "euromoby.com/smsgw/internal/outbound/models"

type TestBilling struct{}

func NewTestBilling() *TestBilling {
	return &TestBilling{}
}

func (b *TestBilling) CheckBalance(merchantID string) error {
	return nil
}

func (b *TestBilling) ChargeOutboundMessage(message *om.Message) error {
	return nil
}
