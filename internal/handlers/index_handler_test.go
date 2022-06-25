package handlers

import (
	"encoding/json"
	"euromoby.com/smsgw/internal/testhelpers"
	"net/http"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	h := NewIndexHandler()

	recorder := testhelpers.PerformAnonRequest(h.Index, http.MethodGet, "/", nil)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}

	var indexResponse IndexResponse
	if err := json.NewDecoder(recorder.Body).Decode(&indexResponse); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if indexResponse.Health != "OK" {
		t.Errorf("App is not healthy")
	}
}
