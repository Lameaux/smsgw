package middlewares

import (
	"errors"
	"net/http"

	"euromoby.com/smsgw/internal/config"
	"github.com/gin-gonic/gin"
)

const (
	HeaderXApiKey = "X-Api-Key"

	// MerchantIDKey is the key that holds the merchant ID in a request context.
	MerchantIDKey = "MerchantID"
)

var ErrUnauthorized = errors.New("Unauthorized")

type Authenticator struct {
	appConfig *config.AppConfig
}

func NewAuthenticator(appConfig *config.AppConfig) *Authenticator {
	return &Authenticator{appConfig}
}

func (auth *Authenticator) Authenticate(c *gin.Context) {
	merchant, err := auth.doAuthenticate(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, err)
		return
	}
	c.Set(MerchantIDKey, merchant)
	c.Next()
}

func (auth *Authenticator) doAuthenticate(r *http.Request) (string, error) {
	merchant, exists := auth.appConfig.Merchants[r.Header.Get(HeaderXApiKey)]

	if !exists {
		return "", ErrUnauthorized
	}

	return merchant, nil
}
