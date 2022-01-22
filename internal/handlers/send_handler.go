package handlers

import (
	"encoding/json"
	"net/http"

	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/services"
	"euromoby.com/smsgw/internal/utils"
	"github.com/gin-gonic/gin"
)

type SendHandler struct {
	service *services.MessageOrderService
}

func NewSendHandler(service *services.MessageOrderService) *SendHandler {
	return &SendHandler{service}
}

func (h *SendHandler) SendMessage(c *gin.Context) {
	p, err := h.parseRequest(c)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.SendMessage(p)
	if err != nil {
		switch err {
		case models.ErrDuplicateClientTransactionID:
			c.JSON(http.StatusConflict, result)
		default:
			utils.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *SendHandler) parseRequest(c *gin.Context) (*inputs.SendMessageParams, error) {
	var p inputs.SendMessageParams
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&p)
	if err != nil {
		return nil, err
	}

	p.MerchantID = c.GetString(middlewares.MerchantIDKey)

	recipients, err := h.normalizeRecipients(p.To)
	if err != nil {
		return nil, err
	}
	p.To = recipients

	// TODO: validate more inputs

	return &p, nil
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
