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
)

func insertWebhookResponseToTable(body string, response string, code int) {
	_, err := InsertIntoTable(&data.WebhookMessage{
		Message:      body,
		Response:     response,
		CodeResponse: code,
		Timestamp:    time.Now(),
	})
	if err != nil {
		failOnError(err, "Error inserting webhook response into table")
	}
}

func SendWebhook(message data.StoredMessage, clientID string, webhookURL string) {
	if webhookURL != "" {
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
		insertWebhookResponseToTable(string(jsonData), responseBody, statusCode)
	}
}
