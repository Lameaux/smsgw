package inputs

type MessageOrderParams struct {
	MerchantID string
	ID         string
}

type MessageOrderSearchParams struct {
	*SearchParams

	MerchantID          string
	ClientTransactionID *string
}
