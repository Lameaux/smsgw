package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/services"
	"euromoby.com/smsgw/internal/views"
)

type InboundHandler struct {
	service *services.InboundService
}

func NewInboundHandler(service *services.InboundService) *InboundHandler {
	return &InboundHandler{service}
}

func (h *InboundHandler) Get(c *gin.Context) {
	p := h.params(c)

	message, err := h.service.FindByMerchantAndID(p.MerchantID, p.ID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			views.ErrorJSON(c, http.StatusNotFound, models.ErrMessageNotFound)
		} else {
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *InboundHandler) Ack(c *gin.Context) {
	p := h.params(c)

	m, err := h.service.AckByMerchantAndID(p.MerchantID, p.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrAlreadyAcked):
			c.JSON(http.StatusConflict, m)
		case errors.Is(err, models.ErrNotFound):
			views.ErrorJSON(c, http.StatusNotFound, models.ErrMessageNotFound)
		default:
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusOK, m)
}

func (h *InboundHandler) Search(c *gin.Context) {
	p, err := h.searchParams(c)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	messages, err := h.service.FindByQuery(p)
	if err != nil {
		views.ErrorJSON(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *InboundHandler) params(c *gin.Context) *inputs.InboundMessageParams {
	return &inputs.InboundMessageParams{
		MerchantID: c.GetString(middlewares.MerchantIDKey),
		ID:         c.Param("id"),
	}
}

func (h *InboundHandler) searchParams(c *gin.Context) (*inputs.InboundMessageSearchParams, error) {
	sp, err := commonSearchParams(c)
	if err != nil {
		return nil, err
	}

	mp, err := messageSearchParams(c)
	if err != nil {
		return nil, err
	}

	p := inputs.InboundMessageSearchParams{
		MerchantID:    c.GetString(middlewares.MerchantIDKey),
		SearchParams:  sp,
		MessageParams: mp,
	}

	if shortcode := c.Query("shortcode"); shortcode != "" {
		p.Shortcode = &shortcode
	}

	return &p, nil
}
