package middleware

import (
	"homework/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("Invalid token claims")
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	username := claims["sub"].(string)
	log.Printf("Token verification successful for user: %s\n. Claims %+v\n", username, token.Claims)
	c.Set("username", username)

	c.Next()
}
