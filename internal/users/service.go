package users

type Service interface {
	Authorize(apiKey string) (string, error)
	FindMerchantByShortcode(shortcode string) (string, error)
}
