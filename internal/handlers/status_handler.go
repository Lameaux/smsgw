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

type StatusHandler struct {
	service *services.MessageOrderService
}

func NewStatusHandler(service *services.MessageOrderService) *StatusHandler {
	return &StatusHandler{service}
}

func (h *StatusHandler) Get(c *gin.Context) {
	p := h.params(c)

	orderStatus, err := h.service.FindByMerchantAndID(p.MerchantID, p.ID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			views.ErrorJSON(c, http.StatusNotFound, models.ErrMessageOrderNotFound)
		} else {
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusOK, orderStatus)
}

func (h *StatusHandler) Search(c *gin.Context) {
	p, err := h.searchParams(c)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	messageOrders, err := h.service.FindByQuery(p)
	if err != nil {
		views.ErrorJSON(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, messageOrders)
}

func (h *StatusHandler) params(c *gin.Context) *inputs.MessageOrderParams {
	return &inputs.MessageOrderParams{
		MerchantID: c.GetString(middlewares.MerchantIDKey),
		ID:         c.Param("id"),
	}
}

func (h *StatusHandler) searchParams(c *gin.Context) (*inputs.MessageOrderSearchParams, error) {
	sp, err := commonSearchParams(c)
	if err != nil {
		return nil, err
	}

	p := inputs.MessageOrderSearchParams{
		MerchantID:   c.GetString(middlewares.MerchantIDKey),
		SearchParams: sp,
	}

	clientTransactionID := c.Query("client_transaction_id")
	if clientTransactionID != "" {
		p.ClientTransactionID = &clientTransactionID
	}

	return &p, nil
}
