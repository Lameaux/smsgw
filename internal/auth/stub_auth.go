package auth

import "euromoby.com/smsgw/internal/models"

type StubAuth struct {
	Merchants map[string]string
}

func NewStubAuth() *StubAuth {
	merchants := map[string]string{
		"postman-api-key": "d70c94da-dac4-4c0c-a6db-97f1740f29a8",
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

func (a *StubAuth) ValidateShortcode(merchantID, shortcode string) error {
	return nil
}
