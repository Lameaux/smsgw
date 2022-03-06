package auth

import "euromoby.com/smsgw/internal/models"

type TestAuth struct {
	Merchants map[string]string
}

const (
	TestAPIKey     = "test-api-key"
	TestMerchantID = "d70c94da-dac4-4c0c-a6db-97f1740f29aa"
)

func NewTestAuth() *TestAuth {
	merchants := map[string]string{
		TestAPIKey: TestMerchantID,
	}

	return &TestAuth{merchants}
}

func (a *TestAuth) Authorize(apiKey string) (string, error) {
	merchant, exists := a.Merchants[apiKey]
	if !exists {
		return "", models.ErrNotFound
	}

	return merchant, nil
}
