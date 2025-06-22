package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"users/handlers"
	"users/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRegisterAndLogin(t *testing.T) {
	// Register
	payload := map[string]string{"email": "test@example.com", "password": "securepass"}
	if c := request(payload, "POST", "/register").Code; c != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", c)
	}

	// Login
	res := request(payload, "POST", "/login")
	if res.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.Code)
	}

	var resp map[string]string
	json.Unmarshal(res.Body.Bytes(), &resp)
	if _, exists := resp["token"]; !exists {
		t.Errorf("expected token in response")
	} else {
		t.Logf("Received token: %s", resp["token"])
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	// Attempt to login with unregistered user
	payload := map[string]string{"email": "abc@gas.com", "password": "wrongpass"}
	
	if c := request(payload, "POST", "/login").Code; c != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", c)
	}
}

func TestLoginInvalidCredentialsExistingEmail(t *testing.T) {
	payload := map[string]string{"email": "bobby@gmail.com", "password": "correctpass"}
	if c := request(payload, "POST", "/login").Code; c != http.StatusCreated {
		t.Errorf("expected 201 Unauthorized, got %d", c)
	}

	wrong_payload := map[string]string{"email": "bobby@gmail.com", "password": "wrongpass"}
	if c := request(wrong_payload, "POST", "/login").Code; c != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", c)
	}

	if c := request(payload, "POST", "/login").Code; c != http.StatusOK {
		t.Errorf("expected 200 Accepted, got %d", c)
	}
}

func TestAlreadyRegisteredEmail(t *testing.T) {
	payload := map[string]string{"email": "cat@abc.com", "password": "mypassword"}
	if c := request(payload, "POST", "/register").Code; c != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", c)
	}

	payload = map[string]string{"email": "cat@abc.com", "password": "anotherpassword"}
	if c := request(payload, "POST", "/register").Code; c != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", c)
	}
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	models.DB, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	models.ConnectDatabase()
	handlers.RegisterAuthRoutes(r, models.DB)
	return r
}

func request(payload map[string]string, method, url string) *httptest.ResponseRecorder {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r := setupTestRouter()
	r.ServeHTTP(w, req)
	return w
}