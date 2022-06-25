package message

import commoninputs "euromoby.com/smsgw/internal/common/inputs"

type Params struct {
	MerchantID string
	ID         string
}

type SearchParams struct {
	*commoninputs.SearchParams
	*commoninputs.MessageParams

	MerchantID string
}
