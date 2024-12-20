package handlers

import (
	"database/sql"
	"homework/app/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func HandleRegistrationEvent(c *gin.Context, db *sqlx.DB) {
	var payload models.RegistrationEventPayload

	if err := c.BindJSON(&payload); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	var event models.Event
	query := "SELECT id FROM events WHERE id = $1"
	log.Printf("Executing query: %s with EventID: %d", query, payload.EventID)
	err := db.Get(&event, query, payload.EventID)
	if err != nil {
		log.Printf("Error fetching event: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event not found"})
		return
	}

	var existingRegistration models.Registration
	query = "SELECT id FROM registrations WHERE event_id = $1 AND participant_id = $2"
	log.Printf("Checking registration with query: %s (event_id: %d, participant_id: %d)", query, payload.EventID, payload.ParticipantID)
	err = db.Get(&existingRegistration, query, payload.EventID, payload.ParticipantID)
	if err == nil {
		log.Printf("Participant already registered (Registration ID: %d)", existingRegistration.ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Participant already registered"})
		return
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking registration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check registration"})
		return
	}

	registration := models.Registration{
		EventID:          payload.EventID,
		ParticipantID:    payload.ParticipantID,
		RegistrationDate: time.Now(),
	}
	query = "INSERT INTO registrations (event_id, participant_id, registration_date) VALUES (:event_id, :participant_id, :registration_date) RETURNING id"
	log.Printf("Inserting registration with query: %s", query)
	rows, err := db.NamedQuery(query, registration)
	if err != nil {
		log.Printf("Error inserting registration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register participant"})
		return
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&registration.ID)
		if err != nil {
			log.Printf("Error scanning registration ID: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve registration ID"})
			return
		}
	}

	query = "UPDATE events SET participant_count = participant_count + 1 WHERE id = $1"
	log.Printf("Updating participant count with query: %s (event_id: %d)", query, payload.EventID)
	_, err = db.Exec(query, payload.EventID)
	if err != nil {
		log.Printf("Error updating participant count: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update participant count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "registration_id": registration.ID})
}

func HandleListRegistrations(c *gin.Context, db *sqlx.DB) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	var user models.User
	query := "select id from users where username = $1"
	err := db.Get(&user, query, username)
	if err != nil {
		log.Printf("error fetching user: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to fetch user"})
		return
	}

	var registrations []models.Registration
	query = "select event_id from registrations where participant_id = $1"

	err = db.Select(&registrations, query, user.ID)
	if err != nil {
		log.Printf("error fetching registrations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch registrations"})
		return
	}

	var events []models.Event
	eventids := make([]int64, len(registrations))
	for i, reg := range registrations {
		eventids[i] = int64(reg.EventID)
	}

	if len(eventids) > 0 {
		query = `select * from events where id = any($1)`
		err = db.Select(&events, query, pq.Array(eventids))
		if err != nil {
			log.Printf("error fetching events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"events": events})

}
