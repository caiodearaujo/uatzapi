package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
	"whatsgoingon/data"
	"whatsgoingon/helpers"
)

const (
	// ErrNoRows represents the common "no rows in result set" error message.
	ErrNoRows = "sql: no rows in result set"
)

// WebhookBody represents the structure of the request body for adding a webhook.
type WebhookBody struct {
	DeviceID   int    `json:"device_id"`
	WebhookURL string `json:"webhook_url"`
}

// WebhookAdd adds a new webhook to a device.
// It requires a valid device ID and a webhook URL in the request body.
func WebhookAdd(c *gin.Context) {
	var body WebhookBody

	// Bind and validate JSON request body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if body.DeviceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}
	if body.WebhookURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "webhook_url is required"})
		return
	}

	// Add webhook using helper function
	if err, _ := helpers.AddWebhook(body.DeviceID, body.WebhookURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// WebhookList returns a list of all webhooks.
// This route retrieves all existing webhooks using the helper function.
func WebhookList(c *gin.Context) {
	webhookList, err := helpers.ListWebhooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhookList)
}

// WebhookListByDevice returns all webhooks for a specific device.
// It requires the device ID as a URL parameter.
func WebhookListByDevice(c *gin.Context) {
	deviceID := c.Param("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}

	// Convert deviceID to integer
	dvcID, err := strconv.Atoi(deviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}

	// Retrieve webhooks by device ID
	webhookList, err := helpers.ListWebhooksByDeviceID(dvcID)
	if err != nil {
		if err.Error() == ErrNoRows {
			// No webhooks found for this device
			c.JSON(http.StatusOK, []data.DeviceWebhook{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhookList)
}

// WebhookByDevice retrieves the active webhook for a specific device by its ID.
// It requires the device ID as a URL parameter.
func WebhookByDevice(c *gin.Context) {
	deviceID := c.Param("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}

	// Convert deviceID to integer
	dvcID, err := strconv.Atoi(deviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}

	// Retrieve active webhook for the device
	webhook, err := helpers.GetWebhookActiveByDeviceID(dvcID)
	if err != nil {
		if err.Error() == ErrNoRows {
			c.JSON(http.StatusNoContent, gin.H{"status": "no webhooks found for the given device ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhook)
}

// WebhookRemove deactivates (removes) the webhook for a specific device by its ID.
// It requires the device ID as a URL parameter.
func WebhookRemove(c *gin.Context) {
	deviceID := c.Param("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}

	// Convert deviceID to integer
	dvcID, err := strconv.Atoi(deviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}

	// Remove webhook using the helper function
	if err, _ := helpers.RemoveWebhook(dvcID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
