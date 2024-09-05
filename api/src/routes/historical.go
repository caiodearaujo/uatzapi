package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HistoricalConversation(c *gin.Context) {
	recipientID := c.Query("recipient_id")
	if recipientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'recipient_id' is required"})
		return
	}

	messages:= "oi"
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	// 	return
	// }
	c.JSON(http.StatusOK, messages)
	return
}
