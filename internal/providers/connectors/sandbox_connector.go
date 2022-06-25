package connectors

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"euromoby.com/core/logger"
	coremodels "euromoby.com/core/models"
	"euromoby.com/smsgw/internal/config"
)

const (
	apiBaseURL = "http://0.0.0.0:8081/sandbox"
)

type SandboxConnector struct {
	app          *config.App
	countryCodes []string
}

type SandboxMessageRequest struct {
	Sender              string `json:"sender"`
	MSISDN              string `json:"msisdn"`
	Body                string `json:"body"`
	ClientTransactionID string `json:"client_transaction_id"`
}

type SandboxMessageResponse struct {
	MessageID *string `json:"message_id"`
}

func NewSandboxConnector(app *config.App) *SandboxConnector {
	return &SandboxConnector{
		app:          app,
		countryCodes: []string{"420", "357", "380"},
	}
}

func (c *SandboxConnector) Name() string {
	return "sandbox"
}

func (c *SandboxConnector) Accept(message *SendMessageRequest) bool {
	recipient := message.MSISDN.String()

	for _, countryCode := range c.countryCodes {
		if strings.HasPrefix(recipient, countryCode) {
			return true
		}
	}

	return false
}

func (c *SandboxConnector) SendMessage(message *SendMessageRequest) (*SendMessageResponse, error) {
	reqBody := SandboxMessageRequest{
		MSISDN:              message.MSISDN.String(),
		Body:                message.Body,
		Sender:              message.Sender,
		ClientTransactionID: message.ClientTransactionID,
	}

	httpResp, err := c.app.Config.HTTPClient.Post(apiBaseURL+"/message", &reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	respBody := string(respBodyBytes)
	r := SendMessageResponse{
		Body: &respBody,
	}

	if httpResp.StatusCode != http.StatusCreated {
		return &r, ErrSendFailed
	}

	var resp SandboxMessageResponse

	if err := json.Unmarshal(respBodyBytes, &resp); err != nil {
		return &r, coremodels.ErrInvalidJSON
	}

	r.MessageID = resp.MessageID

	logger.Infow("message sent via SandboxConnector", "sms", message)

	return &r, nil
}

func (c *SandboxConnector) SendStatus(status *SendStatusRequest) (*SendStatusResponse, error) {
	return nil, ErrSendFailed
}
