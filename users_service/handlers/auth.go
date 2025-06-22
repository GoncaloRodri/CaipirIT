package handlers

import (
	"log"
	"net/http"
	"time"

	"users/middleware"
	"users/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	DB *gorm.DB
}

func (h *AuthHandler) Register(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{Email: creds.Email, Password: string(hashedPassword)}
	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := models.DB.Where("email = ?", creds.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	var memberships []models.Membership
	models.DB.Where("user_id = ?", user.ID).Find(&memberships)

	orgSet := make(map[uint]struct{})
	teamSet := make(map[uint]struct{})
	for _, m := range memberships {
		teamSet[m.TeamID] = struct{}{}
		var t models.Team
		models.DB.Select("organization_id").First(&t, m.TeamID)
		orgSet[t.OrganizationID] = struct{}{}
	}

	orgIDs := make([]uint, 0, len(orgSet))
	for id := range orgSet {
		orgIDs = append(orgIDs, id)
	}
	teamIDs := make([]uint, 0, len(teamSet))
	for id := range teamSet {
		teamIDs = append(teamIDs, id)
	}

	expiresAt := time.Now().Add(time.Hour)
	claims := &middleware.Claims{
		Email:   user.Email,
		OrgIDs:  orgIDs,
		TeamIDs: teamIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB) {
	authHandler := &AuthHandler{DB: db}

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	log.Println("Authentication routes registered!")
}
