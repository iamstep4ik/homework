package middleware

import (
	"homework/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		log.Println("Token missing in cookie")
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	token, err := utils.VerifyToken(tokenString)
	if err != nil {
		log.Printf("Token verification failed: %v\n", err)
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	log.Printf("Token verification successful. Claims %+v\n", token.Claims)

	c.Next()

}
