package sandbox

import (
	"encoding/json"
	"net/http"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/repos"
	"euromoby.com/smsgw/internal/utils"
	"github.com/gin-gonic/gin"
)

type SandboxInboundHandler struct {
	app *config.AppConfig
}

func NewSandboxInboundHandler(app *config.AppConfig) *SandboxInboundHandler {
	return &SandboxInboundHandler{app}
}

func (h *SandboxInboundHandler) ReceiveMessage(c *gin.Context) {
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

	inboundMessageRepo := repos.NewInboundMessageRepo(conn)
	m := h.makeInboundMessage(mreq)

	err = inboundMessageRepo.Save(m)
	if err != nil {
		switch err {
		case models.ErrDuplicateProviderMessageID:
			utils.ErrorJSON(c, http.StatusConflict, err)
		default:
			utils.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusCreated, &m)
}

func (h *SandboxInboundHandler) parseRequest(r *http.Request) (*SandboxInboundMessage, error) {
	var mreq SandboxInboundMessage
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&mreq)
	if err != nil {
		return nil, err
	}

	msisdn, err := utils.NormalizeMSISDN(mreq.MSISDN)
	if err != nil {
		return nil, err
	}
	mreq.MSISDN = msisdn

	return &mreq, nil
}

func (h *SandboxInboundHandler) makeInboundMessage(mreq *SandboxInboundMessage) *models.InboundMessage {
	now := utils.Now()
	return &models.InboundMessage{
		Shortcode:         mreq.Shortcode,
		MSISDN:            mreq.MSISDN,
		Body:              mreq.Body,
		ProviderID:        SandboxProviderID,
		ProviderMessageID: mreq.MessageID,
		Status:            models.InboundMessageStatusNew,
		NextAttemptAt:     now,
		AttemptCounter:    0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
