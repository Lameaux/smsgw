package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/config"
)

const (
	HeaderXApiKey = "X-Api-Key" //nolint:gosec

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

func (a *Authenticator) Authenticate(c *gin.Context) {
	merchant, err := a.doAuthenticate(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, err)

		return
	}

	c.Set(MerchantIDKey, merchant)
	c.Next()
}

func (a *Authenticator) doAuthenticate(r *http.Request) (string, error) {
	merchant, err := a.appConfig.Auth.Authorize(r.Header.Get(HeaderXApiKey))
	if err != nil {
		return "", ErrUnauthorized
	}

	return merchant, nil
}
