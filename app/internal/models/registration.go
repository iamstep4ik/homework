package models

import (
	"time"
)

type Registration struct {
	ID               int       `json:"id" db:"id"`
	EventID          int       `json:"event_id" db:"event_id"`
	ParticipantID    int       `json:"participant_id" db:"participant_id"`
	RegistrationDate time.Time `json:"registration_date" db:"registration_date"`
}

type RegistrationEventPayload struct {
	EventID       int `json:"event_id" binding:"required"`
	ParticipantID int `json:"participant_id" binding:"required"`
}
