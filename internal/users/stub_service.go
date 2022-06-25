package users

import (
	coremodels "euromoby.com/core/models"
)

type StubService struct {
	Merchants map[string]string
}

const (
	PostmanMerchantID = "d70c94da-dac4-4c0c-a6db-97f1740f29a8"
)

func NewStubAuth() *StubService {
	merchants := map[string]string{
		"postman-api-key": PostmanMerchantID,
		"apikey1":         "d70c94da-dac4-4c0c-a6db-97f1740f29a9",
	}

	return &StubService{merchants}
}

func (a *StubService) Authorize(apiKey string) (string, error) {
	merchant, exists := a.Merchants[apiKey]
	if !exists {
		return "", coremodels.ErrNotFound
	}

	return merchant, nil
}

func (a *StubService) FindMerchantByShortcode(shortcode string) (string, error) {
	return PostmanMerchantID, nil
}
