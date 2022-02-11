package notifiers

import (
	"io"
	"net/http"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/models"
)

type OutboundNotifier struct {
	app *config.AppConfig
}

func NewOutboundNotifier(app *config.AppConfig) *OutboundNotifier {
	return &OutboundNotifier{app}
}

func (on *OutboundNotifier) SendNotification(messageOrder *models.MessageOrder, message *models.OutboundMessage) (*SendNotificationResponse, error) {
	httpResp, err := on.app.HTTPClient.Post(*messageOrder.NotificationURL, &message)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	respBody := string(respBodyBytes)

	r := SendNotificationResponse{
		Body: &respBody,
	}

	if !on.Success(httpResp) {
		return &r, models.ErrSendFailed
	}

	logger.Infow("notification sent", "sms", message)

	return &r, nil
}

func (on *OutboundNotifier) Success(response *http.Response) bool {
	return response.StatusCode >= 200 && response.StatusCode < 300
}
