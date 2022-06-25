package inbound

import commoninputs "euromoby.com/smsgw/internal/common/inputs"

type Params struct {
	MerchantID string
	Shortcode  string
	ID         string
}

type SearchParams struct {
	*commoninputs.SearchParams
	*commoninputs.MessageParams

	MerchantID string
	Shortcode  *string
}
