package inputs

type InboundMessageParams struct {
	MerchantID string
	Shortcode  string
	ID         string
}

type InboundMessageSearchParams struct {
	*SearchParams
	*MessageParams

	MerchantID string
	Shortcode  *string
}
