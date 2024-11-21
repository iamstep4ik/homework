package main

import (
	"homework/config"
	"homework/internal/handlers"
	"homework/internal/middleware"
	"homework/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	database := db.Connect(&cfg)
	defer database.Close()

	r := gin.Default()
	r.LoadHTMLGlob("internal/templates/*")

	r.GET("/register", handlers.ServeRegistrationForm)
	r.POST("/register", func(c *gin.Context) {
		handlers.HandleUserRegistration(c, database)
	})
	r.GET("/login", handlers.ServeLoginForm)
	r.POST("/login", func(c *gin.Context) {
		handlers.HandleUserLogin(c, database)
	})

	r.GET("/home", middleware.AuthMiddleware, func(c *gin.Context) {
		handlers.ServeHomePage(c)
	})

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "localhost", false, true)
		c.Redirect(http.StatusSeeOther, "/login")
	})

	r.Run()

}
