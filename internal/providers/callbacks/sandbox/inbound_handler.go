package sandbox

import (
	"encoding/json"
	"errors"
	"github.com/Lameaux/smsgw/internal/inbound"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Lameaux/core/views"
	"github.com/Lameaux/smsgw/internal/inbound/models"

	coremodels "github.com/Lameaux/core/models"
)

type InboundHandler struct {
	service *inbound.Service
}

func NewInboundHandler(service *inbound.Service) *InboundHandler {
	return &InboundHandler{service}
}

func (h *InboundHandler) ReceiveMessage(c *gin.Context) {
	p, err := h.parseRequest(c.Request)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	merchantID, err := h.service.FindMerchantByShortcode(p.Shortcode)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	m, err := h.makeInboundMessage(merchantID, p)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	err = h.service.SaveMessage(m)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateProviderMessageID):
			c.JSON(http.StatusConflict, m)
		default:
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusCreated, m)
}

func (h *InboundHandler) parseRequest(r *http.Request) (*InboundMessage, error) {
	var p InboundMessage

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	if err := dec.Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (h *InboundHandler) makeInboundMessage(merchantID string, im *InboundMessage) (*models.Message, error) {
	now := coremodels.TimeNow()

	normalized, err := coremodels.NormalizeMSISDN(im.MSISDN)
	if err != nil {
		return nil, err
	}

	m := models.Message{
		MerchantID:        merchantID,
		Shortcode:         im.Shortcode,
		MSISDN:            normalized,
		Body:              im.Body,
		ProviderID:        SandboxProviderID,
		ProviderMessageID: im.MessageID,
		Status:            models.MessageStatusNew,
		NextAttemptAt:     now,
		AttemptCounter:    0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	return &m, nil
}
