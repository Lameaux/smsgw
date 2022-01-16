package handlers

import (
	"net/http"

	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/services"
	"euromoby.com/smsgw/internal/utils"
	"github.com/gin-gonic/gin"
)

type InboundHandler struct {
	service *services.InboundService
}

type InboundParams struct {
	MerchantID string
	Shortcode  string
	ID         string
}

func NewInboundHandler(service *services.InboundService) *InboundHandler {
	return &InboundHandler{service}
}

func (h *InboundHandler) params(c *gin.Context) *InboundParams {
	return &InboundParams{
		MerchantID: c.GetString(middlewares.MerchantIDKey),
		Shortcode:  c.Param("shortcode"),
		ID:         c.Param("id"),
	}
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

	c.JSON(http.StatusOK, &message)
}

func (h *InboundHandler) Ack(c *gin.Context) {
	p := h.params(c)

	err := h.service.ValidateShortcode(p.MerchantID, p.Shortcode)
	if err != nil {
		utils.ErrorJSON(c, http.StatusForbidden, err)
		return
	}

	message, err := h.service.AckByShortcodeAndID(p.Shortcode, p.ID)
	if err != nil {
		switch err {
		case models.ErrAlreadyAcked:
			utils.ErrorJSON(c, http.StatusConflict, err)
		default:
			utils.ErrorJSON(c, http.StatusInternalServerError, err)
		}
		return
	}

	if message == nil {
		utils.ErrorJSON(c, http.StatusNotFound, ErrMessageNotFound)
		return
	}

	c.JSON(http.StatusOK, &message)
}

func (h *InboundHandler) Search(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
