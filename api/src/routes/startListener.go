package routes

import (
	"net/http"
	"whatsgoingon/events"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func StartListener(c *gin.Context) {
	deviceID := c.Query("client_id")
	events.AddToListeners(deviceID)
	c.Status(http.StatusOK)
}
