package routes

import (
	"net/http"
	"strconv"
	"whatsgoingon/helpers"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type WebhookBody struct {
	DeviceID   int `json:"device_id"`
	WebhookURL string `json:"webhook_url"`
}

func WebhookAdd(c *gin.Context) {
	var body WebhookBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if body.DeviceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}

	if body.WebhookURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "webhook_url is required"})
		return
	}

	if err, _ := helpers.AddWebhook(body.DeviceID, body.WebhookURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func WebhookList(c *gin.Context) {
	webhookList, err := helpers.ListWebhooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhookList)
}

func WebhookRemove(c *gin.Context) {
	deviceID := c.Param("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}
	dvcID, _ := strconv.Atoi(deviceID)
	if err, _ := helpers.RemoveWebhook(dvcID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}