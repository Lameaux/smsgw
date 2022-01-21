package connectors

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/models"
)

const (
	apiBaseURL = "http://0.0.0.0:8081/sandbox"
)

type SandboxConnector struct {
	name         string
	countryCodes []string
}

type SandboxConnectorRequest struct {
	Sender string `json:"sender"`
	MSISDN string `json:"msisdn"`
	Body   string `json:"body"`
}
type SandboxConnectorResponse struct {
	MessageID *string `json:"message_id"`
}

func NewSandboxConnector() *SandboxConnector {
	return &SandboxConnector{
		name:         "sandbox",
		countryCodes: []string{"420", "357", "380"},
	}
}

func (c *SandboxConnector) Name() string {
	return c.name
}

func (c *SandboxConnector) Accept(message *MessageRequest) bool {
	for _, countryCode := range c.countryCodes {
		if strings.HasPrefix(message.MSISDN, countryCode) {
			return true
		}
	}
	return false
}

func (c *SandboxConnector) SendMessage(message *MessageRequest) (*MessageResponse, error) {
	reqBody := SandboxConnectorRequest{
		MSISDN: message.MSISDN,
		Body:   message.Body,
		Sender: message.Sender,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	httpResp, err := http.Post(apiBaseURL+"/message", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	respBody := string(respBodyBytes)
	r := MessageResponse{
		Body: &respBody,
	}

	if httpResp.StatusCode != 201 {
		return &r, models.ErrSendFailed
	}

	var resp SandboxConnectorResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return &r, err
	}

	r.MessageID = resp.MessageID

	logger.Infow("message sent via SandboxConnector", "message", message)
	return &r, nil
}

func (c *SandboxConnector) SendStatus(status *StatusRequest) (*StatusResponse, error) {
	return nil, models.ErrSendFailed
}
