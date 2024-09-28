package routes

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
	"whatsgoingon/helpers"
	"whatsgoingon/store"

	"github.com/gin-gonic/gin"
)

type MessageRequest struct {
	DeviceID        int    `json:"device_id"`
	RecipientNumber string `json:"recipient_number"`
	Message         string `json:"message"`
}

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

type MessageResponse struct {
	Status          string    `json:"status"`
	Timestamp       time.Time `json:"timestamp"`
	ID              string    `json:"id"`
	DeviceID        int       `json:"device_id"`
	RecipientNumber string    `json:"recipient_number"`
}

func SendMessage(c *gin.Context) {
	var requestBody MessageRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := requestBody.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jid, err := store.GetJIDByDeviceID(requestBody.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "step": "on connect"})
		return
	}
	resp, err := helpers.SendMessage(jid, requestBody.Message, requestBody.RecipientNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "step": "on connect"})
		return
	}

	response := MessageResponse{
		Status:          "Ok",
		Timestamp:       resp.Timestamp,
		ID:              resp.ID,
		DeviceID:        requestBody.DeviceID,
		RecipientNumber: requestBody.RecipientNumber,
	}

	c.JSON(200, response)
}

func SendSticker(c *gin.Context) {
	dvcID, err := strconv.Atoi(c.Query("device_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requestBody := MessageRequest{
		DeviceID:        dvcID,
		RecipientNumber: c.Query("recipient_number"),
	}

	if err := requestBody.Validate(false); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stickerFile, err := c.FormFile("sticker")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stickerContent, err := stickerFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer stickerContent.Close()

	stickerData, err := io.ReadAll(stickerContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jid, err := store.GetJIDByDeviceID(requestBody.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "step": "on connect"})
		return
	}
	resp, err := helpers.SendSticker(jid, stickerData, requestBody.RecipientNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "step": "on connect"})
		return
	}

	response := MessageResponse{
		Status:          "Ok",
		Timestamp:       resp.Timestamp,
		ID:              resp.ID,
		DeviceID:        requestBody.DeviceID,
		RecipientNumber: requestBody.RecipientNumber,
	}

	c.JSON(200, response)

}
