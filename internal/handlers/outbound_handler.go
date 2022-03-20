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

type OutboundHandler struct {
	service *services.OutboundService
}

func NewOutboundHandler(service *services.OutboundService) *OutboundHandler {
	return &OutboundHandler{service}
}

func (h *OutboundHandler) Get(c *gin.Context) {
	p := h.params(c)

	messageDetail, err := h.service.FindByMerchantAndID(p.MerchantID, p.ID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			views.ErrorJSON(c, http.StatusNotFound, models.ErrMessageNotFound)
		} else {
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusOK, messageDetail)
}

func (h *OutboundHandler) Search(c *gin.Context) {
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

func (h *OutboundHandler) params(c *gin.Context) *inputs.OutboundMessageParams {
	return &inputs.OutboundMessageParams{
		MerchantID: c.GetString(middlewares.MerchantIDKey),
		ID:         c.Param("id"),
	}
}

func (h *OutboundHandler) searchParams(c *gin.Context) (*inputs.OutboundMessageSearchParams, error) {
	sp, err := commonSearchParams(c)
	if err != nil {
		return nil, err
	}

	mp, err := messageSearchParams(c)
	if err != nil {
		return nil, err
	}

	p := inputs.OutboundMessageSearchParams{
		MerchantID:    c.GetString(middlewares.MerchantIDKey),
		SearchParams:  sp,
		MessageParams: mp,
	}

	return &p, nil
}
