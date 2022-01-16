package handlers

import (
	"net/http"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/repos"
	"euromoby.com/smsgw/internal/utils"
	"euromoby.com/smsgw/internal/views"
	"github.com/gin-gonic/gin"
)

type OutboundHandler struct {
	app *config.AppConfig
}

func NewOutboundHandler(app *config.AppConfig) *OutboundHandler {
	return &OutboundHandler{app}
}

func (h *OutboundHandler) Get(c *gin.Context) {
	merchantID := c.GetString(middlewares.MerchantIDKey)
	ID := c.Param("id")

	ctx, cancel := repos.DBConnContext()
	defer cancel()

	conn, err := h.app.DBPool.Acquire(ctx)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	defer conn.Release()

	outboundMessageRepo := repos.NewOutboundMessageRepo(conn)

	message, err := outboundMessageRepo.FindByID(merchantID, ID)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	if message == nil {
		utils.ErrorJSON(c, http.StatusNotFound, ErrMessageNotFound)
		return
	}

	messageOrderRepo := repos.NewMessageOrderRepo(conn)

	messageOrder, err := messageOrderRepo.FindByID(merchantID, message.MessageOrderID)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	if messageOrder == nil {
		utils.ErrorJSON(c, http.StatusNotFound, ErrMessageOrderNotFound)
		return
	}

	c.JSON(http.StatusOK, views.NewOutboundMessageDetail(message, messageOrder))
}

func (h *OutboundHandler) Search(c *gin.Context) {
	merchantID := c.GetString(middlewares.MerchantIDKey)

	q, err := makeOutboundMessageQuery(c)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := repos.DBConnContext()
	defer cancel()

	conn, err := h.app.DBPool.Acquire(ctx)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}
	defer conn.Release()

	outboundMessageRepo := repos.NewOutboundMessageRepo(conn)

	messages, err := outboundMessageRepo.FindByQuery(merchantID, q)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, messages)
}

func makeOutboundMessageQuery(c *gin.Context) (*repos.OutboundMessageQuery, error) {
	offset, err := queryParamIntDefault(c, "offset", 0)
	if err != nil {
		return nil, err
	}

	limit, err := queryParamIntDefault(c, "limit", 10)
	if err != nil {
		return nil, err
	}

	createdAtFrom, err := queryParamTime(c, "created_at_from")
	if err != nil {
		return nil, err
	}

	createdAtTo, err := queryParamTime(c, "created_at_to")
	if err != nil {
		return nil, err
	}

	q := repos.OutboundMessageQuery{
		Offset:        offset,
		Limit:         limit,
		CreatedAtFrom: createdAtFrom,
		CreatedAtTo:   createdAtTo,
	}

	msisdn := c.Query("msisdn")
	if msisdn != "" {
		msisdn, err := utils.NormalizeMSISDN(msisdn)
		if err != nil {
			return nil, err
		}

		q.MSISDN = &msisdn
	}

	status := c.Query("status")
	if status != "" {
		q.Status = &status
	}

	return &q, nil
}
