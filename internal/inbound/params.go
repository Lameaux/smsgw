package inbound

import "euromoby.com/smsgw/internal/inputs"

type Params struct {
	MerchantID string
	Shortcode  string
	ID         string
}

type SearchParams struct {
	*inputs.SearchParams
	*inputs.MessageParams

	MerchantID string
	Shortcode  *string
}
