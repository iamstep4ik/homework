package models

import (
	"time"
)

type Event struct {
	ID               int
	Name             string
	Description      string
	Location         string
	StartTime        time.Time
	EndTime          time.Time
	ParticipantCount int
	Date             time.Time
}
