package models

import (
	"time"
)

type Event struct {
	ID               int       `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	Location         string    `json:"location" db:"location"`
	StartTime        time.Time `json:"start_time" db:"start_time"`
	EndTime          time.Time `json:"end_time" db:"end_time"`
	ParticipantCount int       `json:"participant_count" db:"participant_count"`
	Date             time.Time `json:"date_event" db:"date_event"`
	CreatedBy        int       `json:"created_by" db:"created_by"`
}

type CreateEventPayload struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Location    string `json:"location" binding:"required"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
	Date        string `json:"date" binding:"required"`
}
