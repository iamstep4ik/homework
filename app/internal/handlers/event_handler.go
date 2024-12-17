package handlers

import (
	"homework/app/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const timeFormat = "15:04"
const dateFormat = "02-01-06"

func CreateEvent(c *gin.Context, db *sqlx.DB) {
	username := c.MustGet("username").(string)
	var payload models.CreateEventPayload

	if err := c.BindJSON(&payload); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	parsedStartTime, err := time.Parse(timeFormat, payload.StartTime)
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
		return
	}

	parsedEndTime, err := time.Parse(timeFormat, payload.EndTime)
	if err != nil {
		log.Printf("Error parsing end time: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
		return
	}

	parsedDate, err := time.Parse(dateFormat, payload.Date)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	var user models.User

	query := "SELECT id FROM users WHERE username =$1"
	db.Get(&user, query, username)
	event := models.Event{
		Name:             payload.Name,
		Description:      payload.Description,
		Location:         payload.Location,
		StartTime:        parsedStartTime,
		EndTime:          parsedEndTime,
		Date:             parsedDate,
		ParticipantCount: 0,
		CreatedBy:        user.ID,
	}
	query = "INSERT INTO events (name, description, location, start_time,end_time,participant_count,date_event,created_by) VALUES (:name, :description, :location, :start_time,:end_time,:participant_count,:date_event,:created_by) RETURNING id"
	rows, err := db.NamedQuery(query, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		log.Fatalf("Failed to create event %v", err)
		return
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&event.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event id"})
			return
		}
	}
	log.Printf("event successfully created %v\n", event)
	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully", "event_id": event.ID})
}

func HandleMyEvents(c *gin.Context, db *sqlx.DB) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	var user models.User
	query := "SELECT id FROM users WHERE username = $1"
	err := db.Get(&user, query, username)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	var events []models.Event
	query = "SELECT id, name, description, location, start_time, end_time, date_event, participant_count, created_by FROM events WHERE created_by = $1"
	err = db.Select(&events, query, user.ID)
	if err != nil {
		log.Printf("Error fetching events: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)

}
