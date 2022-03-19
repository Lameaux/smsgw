package notifiers

import (
	"io"
	"net/http"

	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/logger"
	"euromoby.com/smsgw/internal/models"
)

type DefaultNotifier struct {
	app *config.AppConfig
}

func NewDefaultNotifier(app *config.AppConfig) *DefaultNotifier {
	return &DefaultNotifier{app}
}

func (dn *DefaultNotifier) SendNotification(notificationURL string, message interface{}) (*NotifierResponse, error) {
	httpResp, err := dn.app.HTTPClient.Post(notificationURL, &message)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	respBody := string(respBodyBytes)

	r := NotifierResponse{
		Body: &respBody,
	}

	if !dn.Success(httpResp) {
		return &r, models.ErrSendFailed
	}

	logger.Infow("notification sent", "sms", message)

	return &r, nil
}

func (dn *DefaultNotifier) Success(response *http.Response) bool {
	return response.StatusCode >= 200 && response.StatusCode < 300
}
