package inputs

type OutboundMessageParams struct {
	MerchantID string
	ID         string
}

type OutboundMessageSearchParams struct {
	*SearchParams
	*MessageParams

	MerchantID string
}
