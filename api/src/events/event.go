package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"whatsgoingon/data"
	"whatsgoingon/handler"
	"whatsgoingon/helpers"
	"whatsgoingon/store"
	myStore "whatsgoingon/store"
)

// Initialize the listener by getting all device IDs and adding them to the listener.
func InitListener() {
	store.BulkUpdateDeviceHandlerOff()

	whatsappIDs, err := helpers.GetAllWhatsappIDs()
	handler.FailOnError(err, "Get deviceIDs failed")

	for _, jid := range whatsappIDs {
		AddToListeners(jid)
	}

}

// Add client to message listeners
func AddToListeners(whatsappID string) {
	_, err := StartMessageListener(whatsappID)
	if err != nil {
		handler.FailOnError(err, fmt.Sprintf("Error starting message listener for %s", whatsappID))
		return
	}

	device, err := store.GetDeviceByJID(whatsappID)
	if err != nil {
		handler.FailOnError(err, fmt.Sprintf("Error getting device by JID: %v", whatsappID))
		return
	}

	// Insert into table if the listener started successfully.
	deviceHandler := &data.DeviceHandler{
		DeviceID: device.ID,
		ActiveAt: time.Now(),
		Active:   true,
	}

	if res, err := store.InsertIntoTable(deviceHandler); err != nil {
		handler.FailOnError(err, fmt.Sprintf("Error inserting into table: %v", res))
	}
}

// Start the message listener for the device ID.
func StartMessageListener(whatsappID string) (*whatsmeow.Client, error) {
	client, err := helpers.GetWhatsAppClientByJID(whatsappID)
	if err != nil {
		handler.FailOnError(err, fmt.Sprintf("Error getting client by ID: %v", whatsappID))
	}

	device, _ := store.GetDeviceByJID(whatsappID)
	webhookUrl, webhookActive, _ := store.GetWebhookURLByDeviceID(device.ID)

	//Add event handler for incoming messages
	client.AddEventHandler(func(evt interface{}) {
		// Handle the message event.
		if msgEvent, ok := evt.(*events.Message); ok {
			handleMessageEvent(msgEvent, client, device.ID, webhookUrl, webhookActive)
		} else if _, ok := evt.(*events.LoggedOut); ok {
			helpers.LogoutDeviceByJID(client.Store.ID.String())
		}
	})

	log.Printf("Starting message listener for %s", device.JID)
	return client, nil
}

// handle the message event, save it to he store, and send async tasks
func handleMessageEvent(msgEvent *events.Message, client *whatsmeow.Client, deviceID int, webhhookURL string, webhookActive bool) {
	ctx := context.Background()
	// Save message to store.
	err, content := myStore.SaveMessage(*msgEvent, client)
	if err != nil {
		handler.FailOnError(err, "Error saving message to store")
		return
	}

	// Process sending to Redis and Webhook concurrently.
	go helpers.SendMessageToRedis(ctx, *content, deviceID)
	go helpers.SendWebhook(*content, deviceID, webhhookURL, webhookActive)
}

// NewClientHandler returns a function that handles the client events.
func NewClientHandler(client *whatsmeow.Client) func(interface{}) {
	return func(evt interface{}) {
		if _, ok := evt.(*events.Connected); ok {
			storeDevice := client.Store
			device := &data.Device{
				JID:          storeDevice.ID.String(),
				PushName:     storeDevice.PushName,
				BusinessName: storeDevice.BusinessName,
				CreatedAt:    time.Now(),
				Active:       true,
			}
			device, err := store.InsertDeviceIfNotExists(device)
			if err != nil {
				handler.FailOnError(err, "Error inserting new device or device already exists")
				return
			}
			AddToListeners(device.JID)

		}
	}
}
