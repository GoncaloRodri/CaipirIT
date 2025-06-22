package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("super_secret_key")

type Claims struct {
	Email 		string 	`json:"email"`
	OrgIDs    	[]uint  `json:"org_ids"`   // orgs the user belongs to
    TeamIDs   	[]uint  `json:"team_ids"`  // teams the user belongs to
    jwt.RegisteredClaims
}

func AuthMiddleware(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
		return
	}
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil || !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
    c.Set("email", claims.Email)
    c.Set("orgIDs", claims.OrgIDs)
    c.Set("teamIDs", claims.TeamIDs)
    c.Next()
}
