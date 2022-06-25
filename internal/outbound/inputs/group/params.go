package group

import "euromoby.com/smsgw/internal/inputs"

type Params struct {
	MerchantID string
	ID         string
}

type SearchParams struct {
	*inputs.SearchParams

	MerchantID          string
	ClientTransactionID *string
}
