package message

import (
	"errors"
	commonhandlers "github.com/Lameaux/smsgw/internal/common/handlers"
	"net/http"

	"github.com/gin-gonic/gin"

	coremodels "github.com/Lameaux/core/models"
	"github.com/Lameaux/core/views"
	"github.com/Lameaux/smsgw/internal/middlewares"
	oim "github.com/Lameaux/smsgw/internal/outbound/inputs/message"
	"github.com/Lameaux/smsgw/internal/outbound/models"
	osm "github.com/Lameaux/smsgw/internal/outbound/services/message"
)

type Handler struct {
	s *osm.Service
}

func NewHandler(s *osm.Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Get(c *gin.Context) {
	p := params(c)

	messageDetail, err := h.s.FindByMerchantAndID(p.MerchantID, p.ID)
	if err != nil {
		if errors.Is(err, coremodels.ErrNotFound) {
			views.ErrorJSON(c, http.StatusNotFound, models.ErrMessageNotFound)
		} else {
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusOK, messageDetail)
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

func params(c *gin.Context) *oim.Params {
	return &oim.Params{
		MerchantID: c.GetString(middlewares.MerchantIDKey),
		ID:         c.Param("id"),
	}
}

func searchParams(c *gin.Context) (*oim.SearchParams, error) {
	sp, err := commonhandlers.CommonSearchParams(c)
	if err != nil {
		return nil, err
	}

	mp, err := commonhandlers.MessageSearchParams(c)
	if err != nil {
		return nil, err
	}

	p := oim.SearchParams{
		MerchantID:    c.GetString(middlewares.MerchantIDKey),
		SearchParams:  sp,
		MessageParams: mp,
	}

	return &p, nil
}
