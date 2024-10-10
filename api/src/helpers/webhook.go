package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"whatsgoingon/data"
	"whatsgoingon/handler"
	"whatsgoingon/store"

	"go.mau.fi/whatsmeow"
)

type WebhookResponse struct {
	Status       string `json:"status"`
	ResponseText string `json:"response_text"`
}

// insertWebhookResponseToTable saves the webhook request/response details in the database.
// It also checks for errors in recent webhook responses to decide if the webhook should be deactivated.
func insertWebhookResponseToTable(deviceID int, webhookURL string, requestBody string, responseBody string, statusCode int) {
	// Insert the webhook response data into the WebhookMessage table.
	_, err := store.InsertIntoTable(&data.WebhookMessage{
		Message:      requestBody,
		Response:     responseBody,
		CodeResponse: statusCode,
		Timestamp:    time.Now(),
		DeviceID:     deviceID,
		WebhookURL:   webhookURL,
	})
	if err != nil {
		handler.FailOnError(err, "Error inserting webhook response into the table")
	}

	// Check if there are frequent errors and deactivate the webhook if necessary.
	inactiveWebhookIfThereAreErrors(deviceID, webhookURL)
}

// inactiveWebhookIfThereAreErrors deactivates the webhook URL if frequent errors are detected within the last 20 messages.
// It performs this check every 5 minutes.
func inactiveWebhookIfThereAreErrors(deviceID int, webhookURL string) {
	// Perform the check every 5 minutes.
	if currentMinute := time.Now().Minute(); currentMinute%5 != 0 {
		webhookMessages, _ := store.GetTop20WebhookMessagesByDeviceID(deviceID)

		// Check if all the recent webhook messages returned non-200 status codes.
		if AllMessagesNon200(webhookMessages) {
			if err := store.InactiveWebhookURLByDeviceID(deviceID); err != nil {
				handler.FailOnError(err, "Error deactivating webhook URL")
				return
			}
			log.Printf("Webhook URL %s deactivated for device ID: %d", webhookURL, deviceID)
		}
	}
}

// AllMessagesNon200 checks if all the provided webhook messages contain non-200 status codes.
func AllMessagesNon200(messages []data.WebhookMessage) bool {
	for _, message := range messages {
		if message.CodeResponse == 200 {
			return false
		}
	}
	return true
}

// SendWebhook sends a webhook message to the specified URL and logs the response in the database.
// It only sends the message if the webhook URL is active.
func SendWebhook(message data.StoredMessage, device data.Device, webhookURL string, webhookActive bool, client *whatsmeow.Client) {
	if webhookURL == "" || !webhookActive {
		return
	}

	// Create a new struct to hold both the message and the device JID.
	message.JID = device.JID

	// Marshal the payload into JSON format.
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal payload to JSON: %v", err)
		return
	}

	// Send the HTTP POST request with the JSON payload.
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to send HTTP request: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return
	}

	var webhookResp WebhookResponse
	err = json.Unmarshal(body, &webhookResp)
	if err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
		return
	}

	// Log the webhook response and status code.
	statusCode := resp.StatusCode
	responseBody := strings.ReplaceAll(string(body), "\n", "")
	insertWebhookResponseToTable(device.ID, webhookURL, string(jsonData), responseBody, statusCode)
	SendMessage(device.JID, webhookResp.ResponseText, message.RecipientID, client)
}

// AddWebhook adds a new webhook for a specific device.
func AddWebhook(deviceID int, webhookURL string) (error, bool) {
	err, _ := store.CreateNewWebhook(deviceID, webhookURL)
	if err != nil {
		return err, false
	}
	return nil, true
}

// RemoveWebhook deactivates a webhook URL for a specific device.
func RemoveWebhook(deviceID int) (error, bool) {
	err := store.InactiveWebhookURLByDeviceID(deviceID)
	if err != nil {
		return fmt.Errorf("error deactivating webhook URL: %v", err), false
	}
	return nil, true
}

// ListWebhooks retrieves all active webhooks from the database.
func ListWebhooks() ([]data.DeviceWebhook, error) {
	return store.GetWebhookURLs()
}

// ListWebhooksByDeviceID retrieves all webhooks associated with a specific device ID.
func ListWebhooksByDeviceID(deviceID int) ([]data.DeviceWebhook, error) {
	return store.GetWebhooksByDeviceID(deviceID)
}

// GetWebhookActiveByDeviceID retrieves the active webhook for a specific device ID.
func GetWebhookActiveByDeviceID(deviceID int) (data.DeviceWebhook, error) {
	return store.GetWebhookActiveByDeviceID(deviceID)
}
