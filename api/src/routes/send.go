package routes

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
	"whatsgoingon/helpers"
	"whatsgoingon/store"

	"github.com/gin-gonic/gin"
)

// MessageRequest represents the request payload for sending messages and stickers.
type MessageRequest struct {
	DeviceID        int    `json:"device_id"`         // ID of the device sending the message
	RecipientNumber string `json:"recipient_number"`  // WhatsApp recipient number
	Message         string `json:"message,omitempty"` // Message text (optional for stickers)
}

// Validate checks if the required fields in MessageRequest are provided.
// Optionally, you can enforce the presence of the Message field with checkMessage argument.
func (m MessageRequest) Validate(checkMessage ...bool) error {
	if m.DeviceID == 0 {
		return errors.New("device_id is required")
	}
	if m.RecipientNumber == "" {
		return errors.New("recipient_number is required")
	}
	if m.Message == "" && len(checkMessage) > 0 && checkMessage[0] {
		return errors.New("message is required")
	}
	return nil
}

// MessageResponse represents the response structure after a message or sticker is sent.
type MessageResponse struct {
	Status          string    `json:"status"`           // Status of the operation
	Timestamp       time.Time `json:"timestamp"`        // Timestamp of the message or sticker sent
	ID              string    `json:"id"`               // WhatsApp message ID
	DeviceID        int       `json:"device_id"`        // ID of the device used to send the message
	RecipientNumber string    `json:"recipient_number"` // WhatsApp recipient number
}

// SendMessage handles the request to send a text message.
// It validates the input, retrieves the device JID, and sends the message via the helper function.
func SendMessage(c *gin.Context) {
	var requestBody MessageRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Validate the request payload
	if err := requestBody.Validate(true); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the JID (WhatsApp ID) based on the device ID
	jid, err := store.GetJIDByDeviceID(requestBody.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve device JID", "details": err.Error()})
		return
	}

	// Send the message using the helper function
	resp, err := helpers.SendMessage(jid, requestBody.Message, requestBody.RecipientNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message", "details": err.Error()})
		return
	}

	// Create and return the response
	response := MessageResponse{
		Status:          "Ok",
		Timestamp:       resp.Timestamp,
		ID:              resp.ID,
		DeviceID:        requestBody.DeviceID,
		RecipientNumber: requestBody.RecipientNumber,
	}
	c.JSON(http.StatusOK, response)
}

// SendSticker handles the request to send a sticker.
// It validates the input, retrieves the sticker file, and sends it via the helper function.
func SendSticker(c *gin.Context) {
	// Parse and validate device ID and recipient number from query parameters
	dvcID, err := strconv.Atoi(c.Query("device_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device_id", "details": err.Error()})
		return
	}
	requestBody := MessageRequest{
		DeviceID:        dvcID,
		RecipientNumber: c.Query("recipient_number"),
	}

	// Validate the request payload (without requiring message)
	if err := requestBody.Validate(false); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the sticker file from the form data
	stickerFile, err := c.FormFile("sticker")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sticker file is required", "details": err.Error()})
		return
	}

	// Open and read the sticker file content
	stickerContent, err := stickerFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open sticker file", "details": err.Error()})
		return
	}
	defer func(stickerContent multipart.File) {
		err := stickerContent.Close()
		if err != nil {
			fmt.Printf("Error closing sticker file: %v\n", err)
		}
	}(stickerContent)

	// Read the entire sticker file content
	stickerData, err := io.ReadAll(stickerContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read sticker file", "details": err.Error()})
		return
	}

	// Retrieve the JID (WhatsApp ID) based on the device ID
	jid, err := store.GetJIDByDeviceID(requestBody.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve device JID", "details": err.Error()})
		return
	}

	// Send the sticker using the helper function
	resp, err := helpers.SendSticker(jid, stickerData, requestBody.RecipientNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send sticker", "details": err.Error()})
		return
	}

	// Create and return the response
	response := MessageResponse{
		Status:          "Ok",
		Timestamp:       resp.Timestamp,
		ID:              resp.ID,
		DeviceID:        requestBody.DeviceID,
		RecipientNumber: requestBody.RecipientNumber,
	}
	c.JSON(http.StatusOK, response)
}
