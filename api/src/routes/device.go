package routes

import (
	"context"
	"net/http"
	"strconv"
	"whatsgoingon/events"
	"whatsgoingon/helpers"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// DeviceList retrieves a list of all devices and returns it as a JSON response.
func DeviceList(c *gin.Context) {
	devices, err := helpers.GetDeviceList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, devices)
}

// DeviceNew creates a new device connection, generates a QR code for client authentication,
// and returns the QR code in base64 format.
func DeviceNew(c *gin.Context) {
	// Initialize new WhatsApp client
	client, err := helpers.NewClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect() // Ensure client is disconnected at the end

	// Get QR code channel for WhatsApp client
	qrChan, _ := client.GetQRChannel(context.Background())

	// Connect client to WhatsApp
	if err := client.Connect(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle QR code events from the channel
	for evt := range qrChan {
		if evt.Event == "code" {
			// Generate QR code in base64 format
			qrCodeBase64, err := helpers.GenerateQRCode(evt.Code)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Add event handler for client
			client.AddEventHandler(events.NewClientHandler(client))

			// Return the QR code to the client
			c.JSON(http.StatusOK, gin.H{"qrCode": qrCodeBase64})
			return
		}
	}
}

// GetDeviceInfo retrieves information about a specific device by its ID and returns it as JSON.
func GetDeviceInfo(c *gin.Context) {
	deviceID := c.Param("deviceId")

	// Convert deviceID from string to integer
	deviceIDInt, err := strconv.Atoi(deviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}

	// Get WhatsApp client by device ID
	client, err := helpers.GetWhatsappClientByDeviceID(deviceIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect() // Ensure client is disconnected at the end

	// Get device information using the client
	deviceInfo := helpers.GetClientInfo(deviceIDInt, client)

	// Return the device information as JSON
	c.JSON(http.StatusOK, deviceInfo)
}
