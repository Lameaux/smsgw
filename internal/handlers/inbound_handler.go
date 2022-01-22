package handlers

import (
	"net/http"

	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/middlewares"
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

func (h *InboundHandler) Get(c *gin.Context) {
	p := h.params(c)

	err := h.service.ValidateShortcode(p.MerchantID, p.Shortcode)
	if err != nil {
		utils.ErrorJSON(c, http.StatusForbidden, err)
		return
	}

	message, err := h.service.FindByShortcodeAndID(p.Shortcode, p.ID)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	if message == nil {
		utils.ErrorJSON(c, http.StatusNotFound, ErrMessageNotFound)
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *InboundHandler) Ack(c *gin.Context) {
	p := h.params(c)

	err := h.service.ValidateShortcode(p.MerchantID, p.Shortcode)
	if err != nil {
		utils.ErrorJSON(c, http.StatusForbidden, err)
		return
	}

	m, err := h.service.AckByShortcodeAndID(p.Shortcode, p.ID)
	if err != nil {
		switch err {
		case models.ErrAlreadyAcked:
			c.JSON(http.StatusConflict, m)
		default:
			utils.ErrorJSON(c, http.StatusInternalServerError, err)
		}
		return
	}

	if m == nil {
		utils.ErrorJSON(c, http.StatusNotFound, ErrMessageNotFound)
		return
	}

	c.JSON(http.StatusOK, m)
}

func (h *InboundHandler) Search(c *gin.Context) {
	p, err := h.searchParams(c)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	err = h.service.ValidateShortcode(p.MerchantID, p.Shortcode)
	if err != nil {
		utils.ErrorJSON(c, http.StatusForbidden, err)
		return
	}

	messages, err := h.service.FindByQuery(p)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *InboundHandler) params(c *gin.Context) *inputs.InboundMessageParams {
	return &inputs.InboundMessageParams{
		MerchantID: c.GetString(middlewares.MerchantIDKey),
		Shortcode:  c.Param("shortcode"),
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
		Shortcode:     c.Param("shortcode"),
		SearchParams:  sp,
		MessageParams: mp,
	}

	return &p, nil
}
