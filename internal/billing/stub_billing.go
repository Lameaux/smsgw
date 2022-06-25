package billing

import (
	"github.com/Lameaux/core/logger"
	om "github.com/Lameaux/smsgw/internal/outbound/models"
)

type StubBilling struct {
	Balances map[string]int64
	Paid     map[string]bool
}

func NewStubBilling() *StubBilling {
	balances := map[string]int64{
		"d70c94da-dac4-4c0c-a6db-97f1740f29a8": 1,
		"d70c94da-dac4-4c0c-a6db-97f1740f29a9": 10, //nolint:gomnd
	}

	paid := map[string]bool{}

	return &StubBilling{balances, paid}
}

func (b *StubBilling) CheckBalance(merchantID string) error {
	logger.Infow("check balance", "merchant", merchantID)

	if b.Balances[merchantID] <= 0 {
		return ErrInsufficientFunds
	}

	return nil
}

func (b *StubBilling) Charge(message *om.Message) error {
	if b.Paid[message.ID] {
		return nil
	}

	logger.Infow("charge outbound message", "message", message)

	balance := b.Balances[message.MerchantID]

	if balance <= 0 {
		return ErrInsufficientFunds
	}

	b.Balances[message.MerchantID] = balance - 1

	b.Paid[message.ID] = true

	return nil
}
