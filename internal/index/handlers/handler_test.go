package handlers

import (
	"encoding/json"
	"euromoby.com/smsgw/internal/testhelpers"
	"net/http"
	"testing"
)

func TestHandler(t *testing.T) {
	h := NewHandler()

	recorder := testhelpers.PerformAnonRequest(h.Index, http.MethodGet, "/", nil)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}

	var response Response
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if response.Health != "OK" {
		t.Errorf("App is not healthy")
	}
}
