package http

import (
	"github.com/Lameaux/smsgw/internal/notifications/notifiers/models"
	"io"
	"net/http"

	"github.com/Lameaux/core/logger"
	"github.com/Lameaux/smsgw/internal/config"
	nm "github.com/Lameaux/smsgw/internal/notifications/models"
)

type Notifier struct {
	app *config.App
}

func NewNotifier(app *config.App) *Notifier {
	return &Notifier{app}
}

func (dn *Notifier) Send(url string, message interface{}) (*models.Response, error) {
	httpResp, err := dn.app.Config.HTTPClient.Post(url, &message)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	respBody := string(respBodyBytes)

	r := models.Response{
		Body: &respBody,
	}

	if !Success(httpResp) {
		return &r, nm.ErrSendFailed
	}

	logger.Infow("notification sent", "sms", message)

	return &r, nil
}

func Success(response *http.Response) bool {
	return response.StatusCode >= 200 && response.StatusCode < 300
}
