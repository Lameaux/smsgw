package sandbox

import (
	"encoding/json"
	"net/http"

	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/services"
	"euromoby.com/smsgw/internal/utils"
	"github.com/gin-gonic/gin"
)

type InboundHandler struct {
	service *services.InboundService
}

func NewInboundHandler(service *services.InboundService) *InboundHandler {
	return &InboundHandler{service}
}

func (h *InboundHandler) ReceiveMessage(c *gin.Context) {
	mreq, err := h.parseRequest(c.Request)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	m := h.makeInboundMessage(mreq)

	err = h.service.SaveMessage(m)
	if err != nil {
		switch err {
		case models.ErrDuplicateProviderMessageID:
			c.JSON(http.StatusConflict, m)
		default:
			utils.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusCreated, m)
}

func (h *InboundHandler) parseRequest(r *http.Request) (*InboundMessage, error) {
	var p InboundMessage
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&p)
	if err != nil {
		return nil, err
	}

	msisdn, err := utils.NormalizeMSISDN(p.MSISDN)
	if err != nil {
		return nil, err
	}
	p.MSISDN = msisdn

	return &p, nil
}

func (h *InboundHandler) makeInboundMessage(mreq *InboundMessage) *models.InboundMessage {
	now := utils.Now()
	return &models.InboundMessage{
		Shortcode:         mreq.Shortcode,
		MSISDN:            mreq.MSISDN,
		Body:              mreq.Body,
		ProviderID:        SandboxProviderID,
		ProviderMessageID: mreq.MessageID,
		Status:            models.InboundMessageStatusNew,
		NextAttemptAt:     now,
		AttemptCounter:    0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
