package billing

import om "github.com/Lameaux/smsgw/internal/outbound/models"

type TestBilling struct{}

func NewTestBilling() *TestBilling {
	return &TestBilling{}
}

func (b *TestBilling) CheckBalance(merchantID string) error {
	return nil
}

func (b *TestBilling) Charge(message *om.Message) error {
	return nil
}
