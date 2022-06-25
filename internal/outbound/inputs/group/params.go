package group

import commoninputs "euromoby.com/smsgw/internal/common/inputs"

type Params struct {
	MerchantID string
	ID         string
}

type SearchParams struct {
	*commoninputs.SearchParams

	MerchantID          string
	ClientTransactionID *string
}
