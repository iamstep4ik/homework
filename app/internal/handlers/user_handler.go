package handlers

import (
	"homework/app/internal/models"
	"homework/app/internal/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func ServeHomePage(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		log.Println("username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the home page", "username": username})
}

func HandleUserRegistration(c *gin.Context, db *sqlx.DB) {
	var payload models.RegistrationPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if !utils.IsValidEmail(payload.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is invalid"})
		return
	}

	user := models.User{
		Username:     payload.Username,
		Email:        payload.Email,
		PasswordHash: string(hashPassword),
		CreatedAt:    time.Now(),
	}

	query := `INSERT INTO users (username, email, password_hash, created_at) VALUES (:username, :email, :password_hash, :created_at) RETURNING id`
	rows, err := db.NamedQuery(query, user)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": user.ID})
}

func HandleUserLogin(c *gin.Context, db *sqlx.DB) {
	var payload models.LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	var user models.User
	var query string

	if strings.Contains(payload.Identifier, "@") {
		query = "SELECT * FROM users WHERE email = $1"
	} else {
		query = "SELECT * FROM users WHERE username = $1"
	}

	err := db.Get(&user, query, payload.Identifier)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or email"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	tokenString, err := utils.CreateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": tokenString})
}

func HandleChangeUsername(c *gin.Context, db *sqlx.DB) {
	username := c.MustGet("username").(string)

	var payload models.ChangeUsernamePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	var user models.User
	query := "SELECT * FROM users WHERE username = $1"
	if err := db.Get(&user, query, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}

	query = "UPDATE users SET username = $1 WHERE username = $2"
	if _, err = tx.Exec(query, payload.NewUsername, username); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	tokenString, err := utils.CreateToken(payload.NewUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Username changed successfully", "new_username": payload.NewUsername})
}

func HandleChangePassword(c *gin.Context, db *sqlx.DB) {
	username := c.MustGet("username").(string)

	var payload models.ChangePasswordPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	var user models.User
	query := "SELECT * FROM users WHERE username = $1"
	if err := db.Get(&user, query, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid current password"})
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	query = "UPDATE users SET password_hash = $1 WHERE username = $2"
	if _, err := db.Exec(query, string(newHash), username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func HandleUserProfile(c *gin.Context, db *sqlx.DB) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	query := "SELECT id, username, email, created_at FROM users WHERE username = $1"
	if err := db.Get(&user, query, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	})
}
