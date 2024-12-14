package main

import (
	"homework/app/internal/config"
	"homework/app/internal/handlers"
	"homework/app/internal/middleware"
	"homework/app/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	database := storage.Connect(&cfg)
	defer database.Close()

	r := gin.Default()
	r.POST("/register", func(c *gin.Context) {
		handlers.HandleUserRegistration(c, database)
	})

	r.POST("/login", func(c *gin.Context) {
		handlers.HandleUserLogin(c, database)
	})

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "localhost", false, true)
		c.JSON(200, gin.H{"message": "Logged out successfully"})
	})

	r.POST("/change-username", middleware.Auth, func(c *gin.Context) {
		handlers.HandleChangeUsername(c, database)
	})

	r.POST("/change-password", middleware.Auth, func(c *gin.Context) {
		handlers.HandleChangePassword(c, database)
	})

	r.GET("/profile", middleware.Auth, func(c *gin.Context) {
		handlers.HandleUserProfile(c, database)
	})

	r.POST("/create-event", middleware.Auth, func(c *gin.Context) {
		handlers.CreateEvent(c, database)
	})

	r.GET("/my-events", middleware.Auth, func(c *gin.Context) {
		handlers.HandleMyEvents(c, database)
	})

	r.POST("/register-event", middleware.Auth, func(c *gin.Context) {
		handlers.HandleRegistrationEvent(c, database)
	})

	r.Run()
}
