package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"euromoby.com/smsgw/internal/config"
)

func TestGin(t *testing.T) {
	app := config.NewTestAppConfig()

	r := Gin(app)

	tests := []struct {
		method     string
		url        string
		statusCode int
	}{
		{"GET", "/dummy", http.StatusNotFound},

		{"GET", "/", http.StatusOK},
		{"GET", "/health", http.StatusOK},

		{"POST", "/v1/sms/messages", http.StatusForbidden},

		{"GET", "/v1/sms/messages/status/1", http.StatusForbidden},
		{"GET", "/v1/sms/messages/status/search", http.StatusForbidden},

		{"GET", "/v1/sms/messages/outbound/1", http.StatusForbidden},
		{"GET", "/v1/sms/messages/outbound/search", http.StatusForbidden},

		{"GET", "/v1/sms/messages/inbound/9999/1", http.StatusForbidden},
		{"GET", "/v1/sms/messages/inbound/9999/search", http.StatusForbidden},

		{"GET", "/v1/sms/callbacks/outbound", http.StatusForbidden},
		{"POST", "/v1/sms/callbacks/outbound", http.StatusForbidden},
		{"DELETE", "/v1/sms/callbacks/outbound/1", http.StatusForbidden},
	}

	for _, tt := range tests {
		request := httptest.NewRequest(tt.method, tt.url, nil)

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		if status := recorder.Code; status != tt.statusCode {
			t.Errorf("Handler returned wrong status code for %s %s. Expected: %d. Got: %d.", tt.method, tt.url, tt.statusCode, status)
		}
	}
}
