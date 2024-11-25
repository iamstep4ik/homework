package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func ServeCreateEventForm(c *gin.Context) {
	c.File("internal/templates/event.html")
}
func HandleEventCreation(c *gin.Context, db *sqlx.DB) {

}
