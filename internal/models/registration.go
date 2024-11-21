package models

import (
	"time"
)

type Registration struct {
	ID               int
	EventID          int
	ParticipantID    int
	RegistrationDate time.Time
	ParticipantCount int
}
