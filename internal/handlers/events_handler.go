package handlers

import (
	"homework/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func ServeCreateEventForm(c *gin.Context) {
	c.File("internal/templates/event.html")
}

const timeFormat = "15:04"
const dateFormat = "2006-01-02"

func CreateEvent(c *gin.Context, db *sqlx.DB) {
	username := c.MustGet("username").(string)
	name := c.PostForm("name")
	description := c.PostForm("description")
	location := c.PostForm("location")
	startTime := c.PostForm("start_time")
	endTime := c.PostForm("end_time")
	date := c.PostForm("date")

	parsedStartTime, err := time.Parse(timeFormat, startTime)
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		c.JSON(400, gin.H{"error": "Invalid start time format"})
		return
	}

	parsedEndTime, err := time.Parse(timeFormat, endTime)
	if err != nil {
		log.Printf("Error parsing end time: %v", err)
		c.JSON(400, gin.H{"error": "Invalid end time format"})
		return
	}

	parsedDate, err := time.Parse(dateFormat, date)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		c.JSON(400, gin.H{"error": "Invalid date format"})
		return
	}

	var user models.User

	query := "SELECT id FROM users WHERE username =$1"
	db.Get(&user, query, username)
	event := models.Event{
		Name:             name,
		Description:      description,
		Location:         location,
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrive event id"})
			return
		}
	}
	log.Printf("event successfully created %v\n", event)
	c.Redirect(http.StatusSeeOther, "/home")

	if err != nil {
		log.Printf("Insert error: %v", err)
		c.JSON(500, gin.H{"error": "Failed to insert event"})
		return
	}

	c.JSON(200, gin.H{"message": "Event created successfully"})
	c.Redirect(http.StatusSeeOther, "/home")
}
