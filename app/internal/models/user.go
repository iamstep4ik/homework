package models

import "time"

type User struct {
	ID              int       `json:"id" db:"id"`
	Username        string    `json:"username" db:"username"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"password_hash" db:"password_hash"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	RegistredEvents *[]Event  `json:"registred_events" db:"registred_events"`
	Token           string    `json:"token"`
}

type RegistrationPayload struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginPayload struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type ChangeUsernamePayload struct {
	NewUsername string `json:"new_username" binding:"required"`
}

type ChangePasswordPayload struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
