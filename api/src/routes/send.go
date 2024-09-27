package routes

import (
	"net/http"
	"time"
	"whatsgoingon/helpers"
	"whatsgoingon/store"

	"github.com/gin-gonic/gin"
)

type MessageResponse struct {
	Status    			string 	  `json:"status"`
	Timestamp 			time.Time `json:"timestamp"`
	ID        			string 	  `json:"id"`
	DeviceID  			string    `json:"device_id"`
	ReicipientNumber 	string    `json:"reicipient_number"`
}

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

	response := MessageResponse{
		Status:    			"Ok",
		Timestamp: 			resp.Timestamp,
		ID:        			resp.ID,
		DeviceID:  			deviceID,
		ReicipientNumber: 	number_reicipient,
	}

	c.JSON(200, response)
}
