package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"caipirIT/handlers"
	"caipirIT/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupEventRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	models.DB, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	models.DB.AutoMigrate(&models.User{}, &models.Event{})

	r.Use(func(c *gin.Context) {
		c.Set("email", "eventuser@example.com")
		c.Next()
	})

	r.POST("/events", handlers.CreateEvent)
	r.GET("/events", handlers.ListEvents)
	return r
}

func createTestUserAndEvent() {
	user := models.User{Email: "eventuser@example.com", Password: "hashedpass"}
	models.DB.Create(&user)
	models.DB.Create(&models.Event{Title: "Initial Event", Org: "DevOrg", Team: "DevTeam", CreatorID: user.ID})
	models.DB.Create(&models.Event{Title: "Second Event", Org: "DevOrg", Team: "DevTeam", CreatorID: user.ID})
}

func TestListEvents(t *testing.T) {
	r := setupEventRouter()
	createTestUserAndEvent()

	req, _ := http.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}
	
	var response []models.Event
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	if len(response) != 2 {
		t.Errorf("Expected 2 events, got %d", len(response))
	}
}

func TestCreateEvent(t *testing.T) {
	r := setupEventRouter()
	createTestUserAndEvent()

	payload := map[string]string{"title": "New Event", "organization": "DevOrg", "team": "DevTeam"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %d", w.Code)
	}
}
