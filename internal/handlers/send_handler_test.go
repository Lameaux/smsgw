package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"euromoby.com/smsgw/internal/auth"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/middlewares"
	"euromoby.com/smsgw/internal/services"
	"euromoby.com/smsgw/internal/views"
)

func TestSendHandler_SendMessage(t *testing.T) {
	service := services.NewMessageOrderService(TestAppConfig)
	h := NewSendHandler(service)

	input := inputs.SendMessageParams{
		Body: "hello world",
		To:   []string{"420123456789"},
	}

	s, err := json.Marshal(input)
	if err != nil {
		t.Errorf("Error creating request body: %v", err)
	}

	recorder := performAuthRequest(h.SendMessage, http.MethodPost, "/", bytes.NewReader(s), header{
		Key:   middlewares.HeaderXApiKey,
		Value: auth.TestAPIKey,
	})

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusCreated, status)
	}

	var output views.MessageOrderDetail
	if err := json.NewDecoder(recorder.Body).Decode(&output); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
}
