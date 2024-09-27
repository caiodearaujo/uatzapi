package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"whatsgoingon/data"
	"whatsgoingon/store"
)

// InsertWebhookResponseToTable inserts the webhook response into the table.
func insertWebhookResponseToTable(deviceID string, webhhookURL string, body string, response string, code int) {
	_, err := store.InsertIntoTable(&data.WebhookMessage{
		Message:      body,
		Response:     response,
		CodeResponse: code,
		Timestamp:    time.Now(),
		DeviceID:     deviceID,
		WebhookURL:   webhhookURL,
	})
	if err != nil {
		FailOnError(err, "Error inserting webhook response into table")
	}
	inactiveWebhookIfThereAreErrors(deviceID, webhhookURL)
}

// InactiveWebhookURLFordeviceID deactivates the webhook URL for the given device ID.
func inactiveWebhookIfThereAreErrors(deviceID string, webhookURL string) {
	if currentMinute := time.Now().Minute(); currentMinute%5 != 0 {
		webhookMessages := store.GetTop20WebhookMessagesByDeviceID(deviceID)
		
		// Get all CodeResponse from webhookMessages and check if there are any errors.
		if ok := AllMessagesNon200(webhookMessages); ok {
			if err := store.InactiveWebhookURLByDeviceID(deviceID); err != nil {
				FailOnError(err, "Error deactivating webhook URL")
				return
			}
			log.Printf("Webhook URL %s deactivated for device ID: %v", webhookURL, deviceID)
		}
	}
}

// AllMessagesNon200 checks if all messages are not 200.
func AllMessagesNon200(messages []data.WebhookMessage) bool {
	for _, message := range messages {
		if message.CodeResponse == 200 {
			return false
		}
	}
	return true
}

// SendWebhook sends the webhook message.
func SendWebhook(message data.StoredMessage, deviceID string, webhookURL string, webhookActive bool) {
	if webhookURL != "" && webhookActive {
		jsonData, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message to JSON: %v", err)
			return
		}

		// Create a new HTTP request with the JSON data.
		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Failed to send HTTP request: %v", err)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Failed to close response body: %v", err)
				return
			}
		}(resp.Body)

		// Read the response body.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			return
		}
		// Insert the webhook response into the table.
		statusCode := resp.StatusCode
		responseBody := strings.ReplaceAll(string(body), "\n", "")
		insertWebhookResponseToTable(deviceID, webhookURL, string(jsonData), responseBody, statusCode)
	}
}
