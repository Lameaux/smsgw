package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/services"
	"euromoby.com/smsgw/internal/views"
)

type SendHandler struct {
	service *services.MessageOrderService
}

const (
	MAX_RECIPIENTS = 50
)

func NewSendHandler(service *services.MessageOrderService) *SendHandler {
	return &SendHandler{service}
}

func (h *SendHandler) SendMessage(c *gin.Context) {
	p, err := h.parseRequest(c)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	result, err := h.service.SendMessage(p)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateClientTransactionID):
			c.JSON(http.StatusConflict, result)
		case errors.Is(err, models.ErrInsufficientFunds):
			views.ErrorJSON(c, http.StatusPaymentRequired, err)
		default:
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *SendHandler) parseRequest(c *gin.Context) (*inputs.SendMessageParams, error) {
	var p inputs.SendMessageParams

	dec := json.NewDecoder(c.Request.Body)

	dec.DisallowUnknownFields()

	if err := dec.Decode(&p); err != nil {
		return nil, err
	}

	p.MerchantID = c.GetString(middlewares.MerchantIDKey)

	if p.Body == "" {
		return nil, models.ErrEmptyBody
	}

	if len(p.To) == 0 {
		return nil, models.ErrMissingRecipients
	}

	recipients, err := h.normalizeRecipients(p.To)
	if err != nil {
		return nil, err
	}

	if len(recipients) > MAX_RECIPIENTS {
		return nil, models.ErrMaxRecipients
	}

	p.Recipients = recipients

	return &p, nil
}

func (h *SendHandler) normalizeRecipients(input []string) ([]models.MSISDN, error) {
	m := make(map[models.MSISDN]struct{})

	for _, msisdn := range input {
		normalized, err := models.NormalizeMSISDN(msisdn)
		if err != nil {
			return nil, err
		}

		m[normalized] = struct{}{}
	}

	output := make([]models.MSISDN, 0, len(m))
	for msisdn := range m {
		output = append(output, msisdn)
	}

	return output, nil
}
