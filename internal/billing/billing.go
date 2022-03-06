package billing

import "euromoby.com/smsgw/internal/models"

type Status struct {
	Active bool
}

type Transaction struct {
	MerchantID string
	ProviderID string
	MSISDN     models.MSISDN
	MessageID  string
}

type ChargeResult struct {
	Success bool
}

type Billing interface {
	GetStatus(merchantID string) *Status
	Charge(transaction Transaction) *ChargeResult
}
