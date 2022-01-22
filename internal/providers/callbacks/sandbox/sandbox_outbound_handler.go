package sandbox

import (
	"encoding/json"
	"net/http"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/repos"
	"euromoby.com/smsgw/internal/utils"
	"github.com/gin-gonic/gin"
)

type OutboundHandler struct {
	app *config.AppConfig
}

func NewOutboundHandler(app *config.AppConfig) *OutboundHandler {
	return &OutboundHandler{app}
}

func (h OutboundHandler) ReceiveStatus(c *gin.Context) {
	mreq, err := h.parseRequest(c.Request)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	ctx, done := repos.DBConnContext()
	defer done()

	conn, err := h.app.DBPool.Acquire(ctx)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	defer conn.Release()

	// TODO: send notification

	c.JSON(http.StatusCreated, &mreq)
}

func (h OutboundHandler) parseRequest(r *http.Request) (*OutboundStatus, error) {
	var mreq OutboundStatus
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&mreq)
	if err != nil {
		return nil, err
	}

	return &mreq, nil
}
