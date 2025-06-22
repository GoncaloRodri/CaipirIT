package test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestOrganizationCreation(t *testing.T) {
	// Create a new organization
	payload := map[string]string{"name": "Test Organization"}
	res := request(payload, "POST", "/organizations")
	if res.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", res.Code)
	}

	var orgResponse map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &orgResponse)
	if orgResponse["name"] != "Test Organization" {
		t.Errorf("expected organization name 'Test Organization', got '%s'", orgResponse["name"])
	}
}