package sandbox

import (
	"encoding/json"
	im "github.com/Lameaux/smsgw/internal/outbound/connectors/inputs/message"
	om "github.com/Lameaux/smsgw/internal/outbound/connectors/outputs/message"
	sim "github.com/Lameaux/smsgw/internal/outbound/connectors/sandbox/inputs/message"
	som "github.com/Lameaux/smsgw/internal/outbound/connectors/sandbox/outputs/message"
	"io"
	"net/http"
	"strings"

	"github.com/Lameaux/smsgw/internal/outbound/connectors/models"

	"github.com/Lameaux/core/logger"
	coremodels "github.com/Lameaux/core/models"
	"github.com/Lameaux/smsgw/internal/config"
)

const (
	apiBaseURL = "http://0.0.0.0:8081/sandbox"
)

type Connector struct {
	app          *config.App
	countryCodes []string
}

func NewConnector(app *config.App) *Connector {
	return &Connector{
		app:          app,
		countryCodes: []string{"420", "357", "380"},
	}
}

func (c *Connector) Name() string {
	return "sandbox"
}

func (c *Connector) Accept(message *im.Request) bool {
	recipient := message.MSISDN.String()

	for _, countryCode := range c.countryCodes {
		if strings.HasPrefix(recipient, countryCode) {
			return true
		}
	}

	return false
}

func (c *Connector) SendMessage(message *im.Request) (*om.Response, error) {
	reqBody := sim.Request{
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
	r := om.Response{
		Body: &respBody,
	}

	if httpResp.StatusCode != http.StatusCreated {
		return &r, models.ErrSendFailed
	}

	var resp som.Response

	if err := json.Unmarshal(respBodyBytes, &resp); err != nil {
		return &r, coremodels.ErrInvalidJSON
	}

	r.MessageID = resp.MessageID

	logger.Infow("message sent via SandboxConnector", "sms", message)

	return &r, nil
}
