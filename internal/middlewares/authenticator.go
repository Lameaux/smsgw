package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/users"
)

const (
	HeaderXApiKey = "X-Api-Key" //nolint:gosec

	// MerchantIDKey is the key that holds the merchant ID in a request context.
	MerchantIDKey = "MerchantID"
)

var ErrUnauthorized = errors.New("Unauthorized")

type Authenticator struct {
	u users.Service
}

func NewAuthenticator(u users.Service) *Authenticator {
	return &Authenticator{u}
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
	merchant, err := a.u.Authorize(r.Header.Get(HeaderXApiKey))
	if err != nil {
		return "", ErrUnauthorized
	}

	return merchant, nil
}
