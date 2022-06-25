package sandbox

import (
	"encoding/json"
	"errors"
	"github.com/Lameaux/smsgw/internal/outbound/listeners/sandbox/inputs"
	om "github.com/Lameaux/smsgw/internal/outbound/models"
	osm "github.com/Lameaux/smsgw/internal/outbound/services/message"
	"net/http"

	"github.com/gin-gonic/gin"

	coremodels "github.com/Lameaux/core/models"
	"github.com/Lameaux/core/views"
)

const (
	SandboxProviderID = "sandbox"
)

type Handler struct {
	service *osm.Service
}

func NewHandler(service *osm.Service) *Handler {
	return &Handler{service}
}

func (h Handler) Ack(c *gin.Context) {
	p, err := parseRequest(c.Request)
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
			views.ErrorJSON(c, http.StatusNotFound, om.ErrMessageNotFound)
		default:
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusOK, m)
}

func parseRequest(r *http.Request) (*inputs.Message, error) {
	var p inputs.Message

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	if err := dec.Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}
