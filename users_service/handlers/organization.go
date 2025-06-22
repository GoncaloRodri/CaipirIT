package handlers

import (
	"log"
	"net/http"
	"strconv"
	"users/models"
	"users/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationHandler struct {
	DB *gorm.DB
}

func (h* OrganizationHandler) createOrganization(c *gin.Context) {
	var org models.Organization
	if err := c.ShouldBindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDIfc, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID := userIDIfc.(uint)

	if err := models.DB.Create(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	orgAdmin := models.OrgAdmin{
		OrganizationID: org.ID,
		UserID: userID,
	}

	if err := h.DB.Create(&orgAdmin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign user as org admin"})
		return
	}

	c.JSON(http.StatusCreated, org.ID)
}

func (h *OrganizationHandler) IsOrgAdmin(c *gin.Context) {
	orgIDParam := c.Param("orgID")
	orgID64, err := strconv.ParseUint(orgIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid orgID"})
		return
	}
	orgID := uint(orgID64)

	userIDIfc, _ := c.Get("userID")
	userID := userIDIfc.(uint)

	var count int64
	h.DB.Model(&models.OrgAdmin{}).
		Where("org_id = ? AND user_id = ? AND role = ?", orgID, userID, "org-admin").
		Count(&count)

	c.JSON(http.StatusOK, gin.H{"isAdmin": count > 0})
}

func (h *OrganizationHandler) getOrganization(c *gin.Context) {
	orgIDParam := c.Param("orgID")
	orgID64, err := strconv.ParseUint(orgIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid orgID"})
		return
	}
	orgID := uint(orgID64)

	userOrgs, _ := c.Get("orgIDs")
	if userOrgs == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not part of any organization"})
		return
	}
	orgIDs := userOrgs.([]uint)


	if !utils.Contains(orgIDs, orgID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not part of this organization"})
		return
	}


	var org models.Organization
	if err := h.DB.Preload("Teams").Preload("Administrators").First(&org, orgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve organization"})
		}
		return
	}
	c.JSON(http.StatusOK, org)
}

func RegisterOrganizationRoutes(r *gin.Engine, db *gorm.DB) {
	handler := &OrganizationHandler{DB: db}

	orgGroup := r.Group("/organizations")
	{
		orgGroup.POST("/", handler.createOrganization)
		orgGroup.GET("/:id", handler.getOrganization)
	}

	log.Println("Organization routes registered!")
}