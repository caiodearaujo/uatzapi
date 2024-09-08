package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"whatsgoingon/events"
)

func StartListener(c *gin.Context) {
	clientID := c.Query("client_id")
	events.AddToListeners(clientID)
	c.Status(http.StatusOK)
}
