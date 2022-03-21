package auth

import "euromoby.com/smsgw/internal/models"

type StubAuth struct {
	Merchants map[string]string
}

const (
	PostmanMerchantID = "d70c94da-dac4-4c0c-a6db-97f1740f29a8"
)

func NewStubAuth() *StubAuth {
	merchants := map[string]string{
		"postman-api-key": PostmanMerchantID,
		"apikey1":         "d70c94da-dac4-4c0c-a6db-97f1740f29a9",
	}

	return &StubAuth{merchants}
}

func (a *StubAuth) Authorize(apiKey string) (string, error) {
	merchant, exists := a.Merchants[apiKey]
	if !exists {
		return "", models.ErrNotFound
	}

	return merchant, nil
}

func (a *StubAuth) FindMerchantByShortcode(shortcode string) (string, error) {
	return PostmanMerchantID, nil
}
