package inputs

type InboundCallbackParams struct {
	MerchantID string
	Shortcode  string `json:"shortcode"`
	URL        string `json:"url"`
}
