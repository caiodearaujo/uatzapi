package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
	"whatsgoingon/data"
	"whatsgoingon/helpers"
	"whatsgoingon/store"
)

type WebhookBody struct {
	DeviceID   int    `json:"device_id"`
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

func WebhookListByDevice(c *gin.Context) {
	deviceID := c.Param("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}
	dvcID, _ := strconv.Atoi(deviceID)
	webhookList, err := helpers.ListWebhooksByDeviceID(dvcID)
	if err != nil {
		if err.Error() == store.NO_ROWS_IN_RESULT_SET {
			c.JSON(http.StatusOK, []data.DeviceWebhook{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhookList)
}

func WebhookByDevice(c *gin.Context) {
	deviceID := c.Param("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}
	dvcID, _ := strconv.Atoi(deviceID)
	webhook, err := helpers.GetWebhookActiveByDeviceID(dvcID)
	if err != nil {
		if err.Error() == store.NO_ROWS_IN_RESULT_SET {
			c.JSON(http.StatusNoContent, gin.H{"status": "no webhooks found for the given device ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhook)
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
