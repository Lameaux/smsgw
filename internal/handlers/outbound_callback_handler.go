package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/repos"
	"euromoby.com/smsgw/internal/views"
)

type OutboundCallbackHandler struct {
	appConfig *config.AppConfig
}

func NewOutboundCallbackHandler(appConfig *config.AppConfig) *OutboundCallbackHandler {
	return &OutboundCallbackHandler{appConfig}
}

func (h *OutboundCallbackHandler) GetCallback(c *gin.Context) {
	p := h.params(c)

	repo := repos.NewOutboundCallbackRepo(h.appConfig.DBPool)

	callback, err := repo.FindByMerchant(p.MerchantID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			views.ErrorJSON(c, http.StatusNotFound, models.ErrCallbackNotFound)
		} else {
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.JSON(http.StatusOK, callback)
}

func (h *OutboundCallbackHandler) RegisterCallback(c *gin.Context) {
	p, err := h.parseRequest(c)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	callback := models.NewSimpleOutboundCallback(p.MerchantID, p.URL)

	repo := repos.NewOutboundCallbackRepo(h.appConfig.DBPool)
	if err := repo.Save(callback); err != nil {
		if errors.Is(err, models.ErrDuplicateCallback) {
			c.JSON(http.StatusConflict, callback)
		} else {
			views.ErrorJSON(c, http.StatusBadRequest, err)
		}

		return
	}

	c.JSON(http.StatusOK, callback)
}

func (h *OutboundCallbackHandler) UpdateCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (h *OutboundCallbackHandler) UnregisterCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (h *OutboundCallbackHandler) params(c *gin.Context) *inputs.OutboundCallbackParams {
	return &inputs.OutboundCallbackParams{
		MerchantID: c.GetString(middlewares.MerchantIDKey),
	}
}

func (h *OutboundCallbackHandler) parseRequest(c *gin.Context) (*inputs.OutboundCallbackParams, error) {
	p := h.params(c)

	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(p); err != nil {
		return nil, err
	}

	return p, nil
}
