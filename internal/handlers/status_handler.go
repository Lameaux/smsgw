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

type StatusHandler struct {
	app *config.AppConfig
}

func NewStatusHandler(app *config.AppConfig) *StatusHandler {
	return &StatusHandler{app}
}

func (h StatusHandler) Get(c *gin.Context) {
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

	messageOrderRepo := repos.NewMessageOrderRepo(conn)

	messageOrder, err := messageOrderRepo.FindByID(merchantID, ID)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	if messageOrder == nil {
		utils.ErrorJSON(c, http.StatusNotFound, ErrMessageOrderNotFound)
		return
	}

	outboundMessageRepo := repos.NewOutboundMessageRepo(conn)

	messages, err := outboundMessageRepo.FindByMessageOrderID(merchantID, messageOrder.ID)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, views.NewMessageOrderStatus(messageOrder, messages))
}

func (h StatusHandler) Search(c *gin.Context) {
	merchantID := c.GetString(middlewares.MerchantIDKey)

	q, err := makeMessageOrderQuery(c)
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

	messageOrderRepo := repos.NewMessageOrderRepo(conn)

	messageOrders, err := messageOrderRepo.FindByQuery(merchantID, q)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, messageOrders)
}

func makeMessageOrderQuery(c *gin.Context) (*repos.MessageOrderQuery, error) {
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

	q := repos.MessageOrderQuery{
		Offset:        offset,
		Limit:         limit,
		CreatedAtFrom: createdAtFrom,
		CreatedAtTo:   createdAtTo,
	}

	clientTransactionID := c.Query("client_transaction_id")
	if clientTransactionID != "" {
		q.ClientTransactionID = &clientTransactionID
	}

	return &q, nil
}
