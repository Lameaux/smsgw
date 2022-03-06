package billing

type StubBilling struct{}

func NewStubBilling() *StubBilling {
	return &StubBilling{}
}

func (b *StubBilling) GetStatus(merchantID string) *Status {
	return &Status{
		Active: true,
	}
}

func (b *StubBilling) Charge(transaction Transaction) *ChargeResult {
	return &ChargeResult{
		Success: true,
	}
}
