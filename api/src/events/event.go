package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"whatsgoingon/data"
	"whatsgoingon/helpers"
	"whatsgoingon/store"
	myStore "whatsgoingon/store"
)

// Initialize the listener by getting all device IDs and adding them to the listener.
func InitListener() {
	store.BulkUpdateDeviceHandlerOff()
	
	whatsappIDs, err := helpers.GetAllWhatsappIDs()
	helpers.FailOnError(err, "Get deviceIDs failed")

	for _, jid := range whatsappIDs {
		AddToListeners(jid)
	}

}
// Add client to message listeners
func AddToListeners(whatsappID string) {
	_, err := StartMessageListener(whatsappID)
	if err != nil {
		helpers.FailOnError(err, fmt.Sprintf("Error starting message listener for %s", whatsappID))
		return
	}

	device, err := store.GetDeviceByJID(whatsappID)
	if err != nil {
		helpers.FailOnError(err, fmt.Sprintf("Error getting device by JID: %v", whatsappID))
		return
	}

	// Insert into table if the listener started successfully.
	deviceHandler := &data.DeviceHandler{
		DeviceID: device.DeviceID(),
		ActiveAt: time.Now(),
		Active: true,
	}

	if res, err := store.InsertIntoTable(deviceHandler); err != nil {
		helpers.FailOnError(err, fmt.Sprintf("Error inserting into table: %v", res))
	}
}

// Start the message listener for the device ID.
func StartMessageListener(whatsappID string) (*whatsmeow.Client, error) {
	client, err := helpers.GetWhatsAppClientByJID(whatsappID)
	if err != nil {
		helpers.FailOnError(err, fmt.Sprintf("Error getting client by ID: %v", whatsappID))
	}

	device, _ := store.GetDeviceByJID(whatsappID)
	webhookUrl, webhookActive, _ := store.GetWebhookURLByDeviceID(device.DeviceID())
	
	//Add event handler for incoming messages
	client.AddEventHandler(func(evt interface{}) {
		log.Printf("2 ---> Event received %T", evt)
		
		// Handle the message event.
		if msgEvent, ok := evt.(*events.Message); ok {
			handleMessageEvent(msgEvent, client, device.DeviceID(), webhookUrl, webhookActive)
		}
	})

	log.Printf("Starting message listener for %s", device.DeviceID())
	return client, nil
}

// handle the message event, save it to he store, and send async tasks
func handleMessageEvent(msgEvent *events.Message, client *whatsmeow.Client, deviceID, webhhookURL string, webhookActive bool) {
	ctx := context.Background()
	// Save message to store.
	err, content := myStore.SaveMessage(*msgEvent, client)
	if err != nil {
		helpers.FailOnError(err, "Error saving message to store")
		return
	}

	// Process sending to Redis and Webhook concurrently.
	go helpers.SendMessageToRedis(ctx, *content, deviceID)
	go helpers.SendWebhook(*content, deviceID, webhhookURL, webhookActive)
}

// NewClientHandler returns a function that handles the client events.
func NewClientHandler(client *whatsmeow.Client) func(interface{}) {
	return func(evt interface{}) {
		log.Printf("1 ---> Event received %T", evt)
		if _, ok := evt.(*events.Connected); ok {
			storeDevice := client.Store
			device := &data.Device{
					JID: 		  storeDevice.ID.String(),
					PushName:     storeDevice.PushName,
					BusinessName: storeDevice.BusinessName,
					CreatedAt:    time.Now(),
					Active: 	  true,
				}
			device, err := store.InsertDeviceIfNotExists(device)
			if err != nil {
				helpers.FailOnError(err, "Error inserting new device or device already exists")
				return
			}
			AddToListeners(device.JID)
			
		}
	}
}