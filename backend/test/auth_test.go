package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"caipirIT/models"
	"caipirIT/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	models.DB, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	models.DB.AutoMigrate(&models.User{}, &models.Event{})
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	return r
}

func TestRegisterAndLogin(t *testing.T) {
	r := setupTestRouter()

	// Register
	payload := map[string]string{"email": "test@example.com", "password": "securepass"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", w.Code)
	}

	// Login
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if _, exists := resp["token"]; !exists {
		t.Errorf("expected token in response")
	}
}
