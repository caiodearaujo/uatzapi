package events

import (
	"fmt"
	"log"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"whatsgoingon/data"
	"whatsgoingon/handler"
	"whatsgoingon/helpers"
	"whatsgoingon/store"
)

// InitListener initializes the message listener for all devices.
// It first deactivates any active device handlers, retrieves all WhatsApp IDs, and adds them to the listener.
func InitListener() {
	// Deactivate any active device handlers before starting.
	store.BulkUpdateDeviceHandlerOff()

	// Get all WhatsApp IDs from the database.
	whatsappIDs, err := helpers.GetAllWhatsappIDs()
	if err != nil {
		handler.FailOnError(err, "Failed to retrieve device IDs")
		return
	}

	// Add each WhatsApp ID to the listeners.
	for _, jid := range whatsappIDs {
		AddToListeners(jid)
	}
}

// AddToListeners adds a device to the message listeners based on its WhatsApp ID.
func AddToListeners(whatsappID string) {
	// Start listening for messages from the device.
	_, err := StartMessageListener(whatsappID)
	if err != nil {
		handler.FailOnError(err, fmt.Sprintf("Error starting message listener for %s", whatsappID))
		return
	}

	// Retrieve the device details using the WhatsApp ID.
	device, err := store.GetDeviceByJID(whatsappID)
	if err != nil {
		handler.FailOnError(err, fmt.Sprintf("Error getting device by JID: %s", whatsappID))
		return
	}

	// Insert a new DeviceHandler into the database to mark the device as active.
	deviceHandler := &data.DeviceHandler{
		DeviceID: device.ID,
		ActiveAt: time.Now(),
		Active:   true,
	}

	if res, err := store.InsertIntoTable(deviceHandler); err != nil {
		handler.FailOnError(err, fmt.Sprintf("Error inserting device handler into table: %v", res))
	}
}

// StartMessageListener starts listening for messages from a specific WhatsApp device.
// It registers a message handler and handles different event types (messages and logout events).
func StartMessageListener(whatsappID string) (*whatsmeow.Client, error) {
	// Get the WhatsApp client for the device.
	client, err := helpers.GetWhatsAppClientByJID(whatsappID)
	if err != nil {
		handler.FailOnError(err, fmt.Sprintf("Error getting client for WhatsApp ID: %s", whatsappID))
		return nil, err
	}

	// Retrieve the device details and webhook information.
	device, _ := store.GetDeviceByJID(whatsappID)
	webhookURL, webhookActive, _ := store.GetWebhookURLByDeviceID(device.ID)

	// Add an event handler for incoming messages and logout events.
	client.AddEventHandler(func(evt interface{}) {
		switch event := evt.(type) {
		case *events.Message:
			handleMessageEvent(event, client, device, webhookURL, webhookActive)
		case *events.LoggedOut:
			helpers.LogoutDeviceByJID(client.Store.ID.String())
		}
	})

	log.Printf("Started message listener for %s", device.JID)
	return client, nil
}

// handleMessageEvent processes incoming messages by storing them and sending them to Redis and Webhook.
// It runs tasks concurrently for performance.
func handleMessageEvent(msgEvent *events.Message, client *whatsmeow.Client, device data.Device, webhookURL string, webhookActive bool) {
	// ctx := context.Background()

	// Convert the event to a stored message format.
	content, err := data.ConvertEventToStoredMessage(*msgEvent, client)
	if err != nil {
		handler.FailOnError(err, "Error saving message to store")
		return
	}

	// Send the message to Redis and Webhook concurrently.
	// go helpers.SendMessageToRedis(ctx, *content, device.ID)
	go helpers.SendWebhook(*content, device, webhookURL, webhookActive, client)
}

// NewClientHandler returns a handler function that is triggered when the client connects.
// It inserts the device into the database if it doesn't exist and starts the message listener.
func NewClientHandler(client *whatsmeow.Client) func(interface{}) {
	return func(evt interface{}) {
		// Check if the event is a successful connection event.
		if _, ok := evt.(*events.Connected); ok {
			// Get the device information from the client store.
			storeDevice := client.Store
			device := &data.Device{
				JID:          storeDevice.ID.String(),
				PushName:     storeDevice.PushName,
				BusinessName: storeDevice.BusinessName,
				CreatedAt:    time.Now(),
				Active:       true,
			}

			// Insert the device into the database if it doesn't exist.
			device, err := store.InsertDeviceIfNotExists(device)
			if err != nil {
				handler.FailOnError(err, "Error inserting new device or device already exists")
				return
			}

			// Add the device to the listeners to start receiving messages.
			AddToListeners(device.JID)
		}
	}
}
