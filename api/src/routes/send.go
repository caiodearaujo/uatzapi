package routes

import (
	"net/http"
	"whatsgoingon/helpers"
	"whatsgoingon/store"

	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	deviceID := c.Query("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "number is required"})
		return
	}
	number_reicipient := c.Query("number_reicipient")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "number is required"})
		return
	}
	message := c.Query("message")
	if message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message is required"})
		return
	}

	jid, err := store.GetJIDByDeviceID(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "step": "on connect"})
		return
	}
	resp, err := helpers.SendMessage(jid, message, number_reicipient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "step": "on connect"})
		return
	}

	c.JSON(200, gin.H{"Status": "Ok", "timestamp": resp.Timestamp, "ID": resp.ID, "ServerID": resp.ServerID})
}
