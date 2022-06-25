package notifiers

import (
	"io"
	"net/http"

	"github.com/Lameaux/core/logger"
	"github.com/Lameaux/smsgw/internal/config"
	nm "github.com/Lameaux/smsgw/internal/notifications/models"
)

type DefaultNotifier struct {
	app *config.App
}

func NewDefaultNotifier(app *config.App) *DefaultNotifier {
	return &DefaultNotifier{app}
}

func (dn *DefaultNotifier) SendNotification(notificationURL string, message interface{}) (*NotifierResponse, error) {
	httpResp, err := dn.app.Config.HTTPClient.Post(notificationURL, &message)
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
		return &r, nm.ErrSendFailed
	}

	logger.Infow("notification sent", "sms", message)

	return &r, nil
}

func (dn *DefaultNotifier) Success(response *http.Response) bool {
	return response.StatusCode >= 200 && response.StatusCode < 300
}
