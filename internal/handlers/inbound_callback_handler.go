package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/auth"
	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/repos"
	"euromoby.com/smsgw/internal/views"
)

type InboundCallbackHandler struct {
	appConfig *config.AppConfig
	auth      auth.Auth
}

func NewInboundCallbackHandler(appConfig *config.AppConfig, auth auth.Auth) *InboundCallbackHandler {
	return &InboundCallbackHandler{appConfig, auth}
}

func (h *InboundCallbackHandler) GetCallback(c *gin.Context) {
	p := h.params(c)

	if err := h.auth.ValidateShortcode(p.MerchantID, p.Shortcode); err != nil {
		views.ErrorJSON(c, http.StatusForbidden, err)

		return
	}

	repo := repos.NewInboundCallbackRepo(h.appConfig.DBPool)

	callback, err := repo.FindByShortcode(p.Shortcode)
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

func (h *InboundCallbackHandler) RegisterCallback(c *gin.Context) {
	p, err := h.parseRequest(c)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	callback := models.NewSimpleInboundCallback(p.Shortcode, p.URL)

	repo := repos.NewInboundCallbackRepo(h.appConfig.DBPool)
	if err := repo.Save(callback); err != nil {
		if errors.Is(err, models.ErrDuplicateCallback) {
			c.JSON(http.StatusConflict, callback)
		} else {
			views.ErrorJSON(c, http.StatusBadRequest, err)
		}

		return
	}

	c.JSON(http.StatusOK, callback)

	c.JSON(http.StatusOK, gin.H{})
}

func (h *InboundCallbackHandler) UpdateCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (h *InboundCallbackHandler) UnregisterCallback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (h *InboundCallbackHandler) params(c *gin.Context) *inputs.InboundCallbackParams {
	return &inputs.InboundCallbackParams{
		MerchantID: c.GetString(middlewares.MerchantIDKey),
		Shortcode:  c.Param("shortcode"),
	}
}

func (h *InboundCallbackHandler) parseRequest(c *gin.Context) (*inputs.InboundCallbackParams, error) {
	p := h.params(c)

	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(p); err != nil {
		return nil, err
	}

	return p, nil
}
