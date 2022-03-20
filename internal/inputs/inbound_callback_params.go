package inputs

type InboundCallbackParams struct {
	MerchantID string
	Shortcode  string
	URL        string `json:"url"`
}
