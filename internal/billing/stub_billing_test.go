package billing_test

import (
	coremodels "github.com/Lameaux/core/models"
	"github.com/Lameaux/smsgw/internal/billing"
	om "github.com/Lameaux/smsgw/internal/outbound/models"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	merchantID = "d70c94da-dac4-4c0c-a6db-97f1740f29a8"
	msisdn     = coremodels.MSISDN(123456789)
)

func TestCheckBalance(t *testing.T) {
	r := require.New(t)

	b := billing.NewStubBilling()

	r.Nil(b.CheckBalance(merchantID), "balance should be positive")

	r.ErrorAs(b.CheckBalance("unknown"), &billing.ErrInsufficientFunds, "balance should be zero")
}

func TestCharge(t *testing.T) {
	r := require.New(t)

	b := billing.NewStubBilling()

	m := om.NewMessage("unknown", "", msisdn)
	m.ID = "message0"
	r.ErrorAs(b.Charge(m), &billing.ErrInsufficientFunds, "balance should be zero")

	m = om.NewMessage(merchantID, "", msisdn)
	m.ID = "message1"
	r.Nil(b.Charge(m), "first charge should be successful")
	r.Nil(b.Charge(m), "same charge twice should be successful")

	m = om.NewMessage(merchantID, "", msisdn)
	m.ID = "message2"
	r.ErrorAs(b.Charge(m), &billing.ErrInsufficientFunds, "next charge should fail")
}
