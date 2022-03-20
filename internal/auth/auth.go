package auth

type Auth interface {
	Authorize(apiKey string) (string, error)
	ValidateShortcode(merchantID, shortcode string) error
}
