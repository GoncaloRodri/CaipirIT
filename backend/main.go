package main

import (
	"github.com/gin-gonic/gin"
	"caipirIT/handlers"
	"caipirIT/middleware"
	"caipirIT/models"
)

func main() {
	r := gin.Default()
	models.ConnectDatabase()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware)
	auth.GET("/events", handlers.ListEvents)
	auth.POST("/events", handlers.CreateEvent)

	r.Run(":8080")
}