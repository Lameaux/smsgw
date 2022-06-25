package sandbox

import (
	"encoding/json"
	"errors"
	om "euromoby.com/smsgw/internal/outbound/models"
	osm "euromoby.com/smsgw/internal/outbound/services/message"
	"net/http"

	"github.com/gin-gonic/gin"

	coremodels "euromoby.com/core/models"
	"euromoby.com/core/views"
)

type OutboundHandler struct {
	service *osm.Service
}

func NewOutboundHandler(service *osm.Service) *OutboundHandler {
	return &OutboundHandler{service}
}

func (h OutboundHandler) Ack(c *gin.Context) {
	p, err := h.parseRequest(c.Request)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	m, err := h.service.AckByProviderAndMessageID(SandboxProviderID, p.MessageID)
	if err != nil {
		switch {
		case errors.Is(err, om.ErrAlreadyAcked):
			c.JSON(http.StatusConflict, m)
		case errors.Is(err, coremodels.ErrNotFound):
			views.ErrorJSON(c, http.StatusNotFound, ErrMessageNotFound)
		default:
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusOK, m)
}

func (h OutboundHandler) parseRequest(r *http.Request) (*OutboundDelivery, error) {
	var mreq OutboundDelivery

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	if err := dec.Decode(&mreq); err != nil {
		return nil, err
	}

	return &mreq, nil
}
