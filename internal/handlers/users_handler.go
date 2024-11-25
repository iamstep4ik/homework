package handlers

import (
	"homework/internal/models"
	"homework/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func ServeRegistrationForm(c *gin.Context) {
	c.File("internal/templates/registration.html")
}

func ServeLoginForm(c *gin.Context) {
	c.File("internal/templates/login.html")
}

func ServeHomePage(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		log.Println("username not found in context")
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	c.HTML(http.StatusOK, "home.html", gin.H{
		"username": username,
	})
}

func ServeChangeUsernameForm(c *gin.Context) {
	c.HTML(http.StatusOK, "change-username.html", gin.H{
		"username": c.MustGet("username").(string),
	})
}

func ServeChangePasswordForm(c *gin.Context) {
	c.File("internal/templates/change-password.html")
}

func HandleUserRegistration(c *gin.Context, db *sqlx.DB) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	check := utils.IsValidEmail(email)
	if !check {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email is invalid"})
		log.Fatalf("email is invalid")
	}

	user := models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashPassword),
		CreatedAt:    time.Now(),
	}

	query := `INSERT INTO users (username, email, password_hash, created_at) VALUES (:username, :email, :password_hash, :created_at) RETURNING id`
	rows, err := db.NamedQuery(query, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		log.Fatalf("Failed to register user: %v", err)
		return
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
			return
		}
	}

	log.Printf("user registred: %v\n", user)
	c.Redirect(http.StatusSeeOther, "/login")
}

func HandleUserLogin(c *gin.Context, db *sqlx.DB) {
	identifier := c.PostForm("identifier")
	password := c.PostForm("password")

	var user models.User
	var query string

	if strings.Contains(identifier, "@") {
		query = "SELECT * FROM users WHERE email = $1"
	} else {
		query = "SELECT * FROM users WHERE username =$1"
	}

	err := db.Get(&user, query, identifier)
	if err != nil {
		log.Printf("Failed to retrieve user: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or email"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}
	loggedInUser := user.Username

	tokenString, err := utils.CreateToken(loggedInUser)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating token")
		return
	}

	log.Printf("Token created %s\n", tokenString)
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/home")

}

func HandleChangeUsername(c *gin.Context, db *sqlx.DB) {
	username := c.MustGet("username").(string)
	newUsername := c.PostForm("newUsername")

	var user models.User

	query := "SELECT * FROM users WHERE username = $1"

	err := db.Get(&user, query, username)
	if err != nil {
		log.Printf("Failed to retrieve user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}

	query = "UPDATE users SET username = $1 WHERE username = $2"
	_, err = tx.Exec(query, newUsername, username)
	if err != nil {
		log.Printf("Failed to update username: %v", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	tokenString, err := utils.CreateToken(newUsername)
	if err != nil {
		log.Printf("Failed to create token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

	c.Set("username", newUsername)
	c.Redirect(http.StatusSeeOther, "/home")
}

func HandleChangePassword(c *gin.Context, db *sqlx.DB) {
	username := c.MustGet("username").(string)
	password := c.PostForm("password")
	newPassword := c.PostForm("newPassword")

	var user models.User
	query := "SELECT password_hash FROM users WHERE username = $1"
	err := db.Get(&user, query, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid password"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	query = "UPDATE users SET password_hash = $1 WHERE username = $2"
	_, err = db.Exec(query, hashedPassword, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/home")
}
