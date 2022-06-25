package send

import (
	"bytes"
	"encoding/json"
	"euromoby.com/smsgw/internal/billing"
	"euromoby.com/smsgw/internal/config"
	ois "euromoby.com/smsgw/internal/outbound/inputs/send"
	"euromoby.com/smsgw/internal/outbound/outputs"
	oss "euromoby.com/smsgw/internal/outbound/services/send"
	"euromoby.com/smsgw/internal/testhelpers"
	"euromoby.com/smsgw/internal/users"
	"net/http"
	"testing"

	"euromoby.com/smsgw/internal/middlewares"
)

func TestSendHandler_SendMessage(t *testing.T) {
	app := config.NewTestApp()
	defer app.Config.Shutdown()

	service := oss.NewService(app, billing.NewTestBilling())
	h := NewHandler(service)

	input := ois.Params{
		Body: "hello world",
		To:   []string{"420123456789"},
	}

	s, err := json.Marshal(input)
	if err != nil {
		t.Errorf("Error creating request body: %v", err)
	}

	recorder := testhelpers.PerformAuthRequest(h.SendMessage, http.MethodPost, "/", bytes.NewReader(s), testhelpers.Header{
		Key:   middlewares.HeaderXApiKey,
		Value: users.TestAPIKey,
	})

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusCreated, status)
	}

	var output outputs.GroupView
	if err := json.NewDecoder(recorder.Body).Decode(&output); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
}
