package group

import (
	"errors"
	commonhandlers "euromoby.com/smsgw/internal/common/handlers"
	"euromoby.com/smsgw/internal/outbound/models"
	"net/http"

	"github.com/gin-gonic/gin"

	coremodels "euromoby.com/core/models"
	"euromoby.com/core/views"
	"euromoby.com/smsgw/internal/middlewares"
	oig "euromoby.com/smsgw/internal/outbound/inputs/group"
	osg "euromoby.com/smsgw/internal/outbound/services/group"
)

type Handler struct {
	s *osg.Service
}

func NewHandler(s *osg.Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Get(c *gin.Context) {
	p := params(c)

	statusView, err := h.s.FindByMerchantAndID(p.MerchantID, p.ID)
	if err != nil {
		if errors.Is(err, coremodels.ErrNotFound) {
			views.ErrorJSON(c, http.StatusNotFound, models.ErrMessageGroupNotFound)
		} else {
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusOK, statusView)
}

func (h *Handler) Search(c *gin.Context) {
	p, err := searchParams(c)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	messageGroups, err := h.s.FindByQuery(p)
	if err != nil {
		views.ErrorJSON(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, messageGroups)
}

func params(c *gin.Context) *oig.Params {
	return &oig.Params{
		MerchantID: c.GetString(middlewares.MerchantIDKey),
		ID:         c.Param("id"),
	}
}

func searchParams(c *gin.Context) (*oig.SearchParams, error) {
	sp, err := commonhandlers.CommonSearchParams(c)
	if err != nil {
		return nil, err
	}

	p := oig.SearchParams{
		MerchantID:   c.GetString(middlewares.MerchantIDKey),
		SearchParams: sp,
	}

	clientTransactionID := c.Query("client_transaction_id")
	if clientTransactionID != "" {
		p.ClientTransactionID = &clientTransactionID
	}

	return &p, nil
}
