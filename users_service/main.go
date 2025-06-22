package main

import (
	"users/handlers"
	"users/models"

	"github.com/gin-gonic/gin"
)


func main() {
	r := gin.Default()
	models.ConnectDatabase()

	handlers.RegisterAuthRoutes(r, models.DB)
	handlers.RegisterOrganizationRoutes(r, models.DB)

	r.Run(":8080")
}