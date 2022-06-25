package inbound

import (
	"errors"
	commonhandlers "github.com/Lameaux/smsgw/internal/common/handlers"
	"net/http"

	"github.com/gin-gonic/gin"

	coremodels "github.com/Lameaux/core/models"
	"github.com/Lameaux/core/views"
	"github.com/Lameaux/smsgw/internal/inbound/models"
	"github.com/Lameaux/smsgw/internal/middlewares"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Get(c *gin.Context) {
	p := params(c)

	message, err := h.s.FindByMerchantAndID(p.MerchantID, p.ID)
	if err != nil {
		if errors.Is(err, coremodels.ErrNotFound) {
			views.ErrorJSON(c, http.StatusNotFound, models.ErrMessageNotFound)
		} else {
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *Handler) Ack(c *gin.Context) {
	p := params(c)

	m, err := h.s.AckByMerchantAndID(p.MerchantID, p.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrAlreadyAcked):
			c.JSON(http.StatusConflict, m)
		case errors.Is(err, coremodels.ErrNotFound):
			views.ErrorJSON(c, http.StatusNotFound, models.ErrMessageNotFound)
		default:
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusOK, m)
}

func (h *Handler) Search(c *gin.Context) {
	p, err := searchParams(c)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	messages, err := h.s.FindByQuery(p)
	if err != nil {
		views.ErrorJSON(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, messages)
}

func params(c *gin.Context) *Params {
	return &Params{
		MerchantID: c.GetString(middlewares.MerchantIDKey),
		ID:         c.Param("id"),
	}
}

func searchParams(c *gin.Context) (*SearchParams, error) {
	sp, err := commonhandlers.CommonSearchParams(c)
	if err != nil {
		return nil, err
	}

	mp, err := commonhandlers.MessageSearchParams(c)
	if err != nil {
		return nil, err
	}

	p := SearchParams{
		MerchantID:    c.GetString(middlewares.MerchantIDKey),
		SearchParams:  sp,
		MessageParams: mp,
	}

	if shortcode := c.Query("shortcode"); shortcode != "" {
		p.Shortcode = &shortcode
	}

	return &p, nil
}
