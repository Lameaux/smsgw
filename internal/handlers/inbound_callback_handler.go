package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/repos"
	"euromoby.com/smsgw/internal/views"
)

type InboundCallbackHandler struct {
	app *config.AppConfig
}

func NewInboundCallbackHandler(app *config.AppConfig) *InboundCallbackHandler {
	return &InboundCallbackHandler{app}
}

func (h *InboundCallbackHandler) ListCallbacks(c *gin.Context) {
	p := h.params(c)

	repo := repos.NewInboundCallbackRepo(h.app.DBPool)

	callbacks, err := repo.FindByMerchant(p.MerchantID)
	if err != nil {
		views.ErrorJSON(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, callbacks)
}

func (h *InboundCallbackHandler) RegisterCallback(c *gin.Context) {
	p, err := h.parseRequest(c)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	callback := models.NewSimpleInboundCallback(p.MerchantID, p.Shortcode, p.URL)
	h.doSaveCallback(c, callback)
}

func (h *InboundCallbackHandler) doSaveCallback(c *gin.Context, callback *models.InboundCallback) {
	repo := repos.NewInboundCallbackRepo(h.app.DBPool)
	if err := repo.Save(callback); err != nil {
		if errors.Is(err, models.ErrDuplicateCallback) {
			c.JSON(http.StatusConflict, callback)
		} else {
			views.ErrorJSON(c, http.StatusBadRequest, err)
		}

		return
	}

	c.JSON(http.StatusCreated, callback)
}

func (h *InboundCallbackHandler) UpdateCallback(c *gin.Context) {
	p, err := h.parseRequest(c)
	if err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	repo := repos.NewInboundCallbackRepo(h.app.DBPool)

	callback, err := repo.FindByMerchantAndShortcode(p.MerchantID, p.Shortcode)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			callback = models.NewSimpleInboundCallback(p.MerchantID, p.Shortcode, p.URL)
			h.doSaveCallback(c, callback)
		} else {
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	callback.URL = p.URL

	if err := repo.Update(callback); err != nil {
		views.ErrorJSON(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, callback)
}

func (h *InboundCallbackHandler) UnregisterCallback(c *gin.Context) {
	p := h.params(c)

	repo := repos.NewInboundCallbackRepo(h.app.DBPool)

	callback, err := repo.FindByMerchantAndShortcode(p.MerchantID, p.Shortcode)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			views.ErrorJSON(c, http.StatusNotFound, models.ErrCallbackNotFound)
		} else {
			views.ErrorJSON(c, http.StatusInternalServerError, err)
		}

		return
	}

	if err := repo.Delete(callback); err != nil {
		views.ErrorJSON(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusNoContent, struct{}{})
}

func (h *InboundCallbackHandler) params(c *gin.Context) *inputs.InboundCallbackParams {
	return &inputs.InboundCallbackParams{
		MerchantID: c.GetString(middlewares.MerchantIDKey),
	}
}

func (h *InboundCallbackHandler) parseRequest(c *gin.Context) (*inputs.InboundCallbackParams, error) {
	p := h.params(c)

	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(p); err != nil {
		return nil, err
	}

	if p.Shortcode == "" {
		return nil, models.ErrInvalidShortcode
	}

	_, err := url.ParseRequestURI(p.URL)
	if err != nil {
		return nil, err
	}

	return p, nil
}
