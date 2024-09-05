package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"whatsgoingon/helpers"
)

func SendMessage(c *gin.Context) {
	number := c.Query("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "number is required"})
		return
	}
	message := c.Query("message")
	if message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message is required"})
		return
	}

	err := helpers.SendMessage(number, message)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "step": "on connect"})
		return
	}

	c.JSON(200, gin.H{"Status": "Ok"})
	return
}
