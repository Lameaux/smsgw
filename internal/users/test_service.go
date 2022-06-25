package users

import coremodels "euromoby.com/core/models"

type TestService struct {
	Merchants map[string]string
}

const (
	TestAPIKey     = "test-api-key"
	TestMerchantID = "d70c94da-dac4-4c0c-a6db-97f1740f29aa"
)

func NewTestAuth() *TestService {
	merchants := map[string]string{
		TestAPIKey: TestMerchantID,
	}

	return &TestService{merchants}
}

func (a *TestService) Authorize(apiKey string) (string, error) {
	merchant, exists := a.Merchants[apiKey]
	if !exists {
		return "", coremodels.ErrNotFound
	}

	return merchant, nil
}

func (a *TestService) FindMerchantByShortcode(shortcode string) (string, error) {
	return TestMerchantID, nil
}
