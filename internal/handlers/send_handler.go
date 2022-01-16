package handlers

import (
	"encoding/json"
	"net/http"

	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/services"
	"euromoby.com/smsgw/internal/utils"
	"euromoby.com/smsgw/internal/views"
	"github.com/gin-gonic/gin"
)

type SendHandler struct {
	service *services.OutboundService
}

func NewSendHandler(service *services.OutboundService) *SendHandler {
	return &SendHandler{service}
}

func (h *SendHandler) SendMessage(c *gin.Context) {
	merchantID := c.GetString(middlewares.MerchantIDKey)

	params, err := h.parseRequest(c.Request)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.SendMessage(merchantID, params)
	if err != nil {
		switch err {
		case models.ErrDuplicateClientTransactionID:
			utils.ErrorJSON(c, http.StatusConflict, err)
		default:
			utils.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *SendHandler) parseRequest(r *http.Request) (*views.SendMessageParams, error) {
	var mreq views.SendMessageParams
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&mreq)
	if err != nil {
		return nil, err
	}

	recipients, err := h.normalizeRecipients(mreq.To)
	if err != nil {
		return nil, err
	}
	mreq.To = recipients

	// TODO: validate more inputs

	return &mreq, nil
}

func (h *SendHandler) normalizeRecipients(input []string) ([]string, error) {
	output := []string{}

	for _, msisdn := range input {
		msisdn, err := utils.NormalizeMSISDN(msisdn)
		if err != nil {
			return nil, err
		}
		output = append(output, msisdn)
	}

	return output, nil
}
