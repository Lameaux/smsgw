package send

import (
	"encoding/json"
	"errors"
	"github.com/Lameaux/smsgw/internal/billing"
	"github.com/Lameaux/smsgw/internal/outbound/models"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Lameaux/core/views"
	"github.com/Lameaux/smsgw/internal/middlewares"
	ois "github.com/Lameaux/smsgw/internal/outbound/inputs/send"
	oss "github.com/Lameaux/smsgw/internal/outbound/services/send"

	coremodels "github.com/Lameaux/core/models"
)

type Handler struct {
	s *oss.Service
}

const (
	maxRecipients = 50
)

func NewHandler(s *oss.Service) *Handler {
	return &Handler{s}
}

func (h *Handler) SendMessage(c *gin.Context) {
	p, err := parseRequest(c)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	result, err := h.s.SendMessage(p)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateClientTransactionID):
			c.JSON(http.StatusConflict, result)
		case errors.Is(err, billing.ErrInsufficientFunds):
			views.ErrorJSON(c, http.StatusPaymentRequired, err)
		default:
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusCreated, result)
}

func parseRequest(c *gin.Context) (*ois.Params, error) {
	var p ois.Params

	dec := json.NewDecoder(c.Request.Body)

	dec.DisallowUnknownFields()

	if err := dec.Decode(&p); err != nil {
		return nil, err
	}

	p.MerchantID = c.GetString(middlewares.MerchantIDKey)

	if p.Body == "" {
		return nil, coremodels.ErrEmptyBody
	}

	if len(p.To) == 0 {
		return nil, models.ErrMissingRecipients
	}

	recipients, err := normalizeRecipients(p.To)
	if err != nil {
		return nil, err
	}

	if len(recipients) > maxRecipients {
		return nil, models.ErrMaxRecipients
	}

	p.Recipients = recipients

	return &p, nil
}

func normalizeRecipients(input []string) ([]coremodels.MSISDN, error) {
	m := make(map[coremodels.MSISDN]struct{})

	for _, msisdn := range input {
		normalized, err := coremodels.NormalizeMSISDN(msisdn)
		if err != nil {
			return nil, err
		}

		m[normalized] = struct{}{}
	}

	output := make([]coremodels.MSISDN, 0, len(m))
	for msisdn := range m {
		output = append(output, msisdn)
	}

	return output, nil
}
