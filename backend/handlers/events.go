// handlers/events.go
package handlers

import (
	"net/http"

	"caipirIT/models"
	"github.com/gin-gonic/gin"
)

func ListEvents(c *gin.Context) {
	email, _ := c.Get("email")
	var user models.User
	if err := models.DB.Where("email = ?", email.(string)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var events []models.Event
	if err := models.DB.Where("creator_id = ?", user.ID).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func CreateEvent(c *gin.Context) {
	email, _ := c.Get("email")
	var user models.User
	if err := models.DB.Where("email = ?", email.(string)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		Title string `json:"title" binding:"required"`
		Org   string `json:"organization"`
		Team  string `json:"team"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := models.Event{
		Title:     input.Title,
		Org:       input.Org,
		Team:      input.Team,
		CreatorID: user.ID,
	}

	if err := models.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}
