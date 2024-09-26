package events

import (
	"log"
	"time"
	"whatsgoingon/data"
	"whatsgoingon/helpers"
	myStore "whatsgoingon/store"

	"github.com/ztrue/tracerr"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func InitListener() {
	helpers.BulkUpdateDeviceHandlerOff()
	clientIds, err := helpers.GetAllClientIDs()
	failOnError(err, "Get clientIds failed")

	for _, clientId := range clientIds {
		AddToListeners(clientId)
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func AddToListeners(clientID string) {
	err, client := StartMessageListener(clientID)
	if err != nil {
		log.Printf("Error starting message listener for %s: %v", clientID, err)
	} else {
		res, err := helpers.InsertIntoTable(&data.DeviceHandler{
			Device: data.Device{
				DeviceID:     clientID,
				PushName:     client.Store.PushName,
				BusinessName: client.Store.BusinessName,
				Timestamp:    time.Now(),
			},
			Active: true,
		})
		if err != nil {
			log.Printf("Error inserting into table: %v", err)
		} else {
			log.Printf("Inserted into table: %v", res)
		}
	}
}

func StartMessageListener(clientID string) (error, *whatsmeow.Client) {
	client, err := helpers.GetClientById(clientID)
	if err != nil {
		tracerr.Print(err)
		return err, nil
	} else {
		webhookUrl, _ := helpers.GetWebhookURLForClientID(clientID)
		client.AddEventHandler(func(evt interface{}) {
			log.Printf("Event received %T", evt)
			switch v := evt.(type) {
			case *events.Message:
				if err, content := myStore.SaveMessage(*v, client); err != nil {
					tracerr.Print(err)
				} else {
					helpers.SendMessageToRedis(*content, clientID)      // async
					helpers.SendWebhook(*content, clientID, webhookUrl) // async
				}
			}
		})
		log.Printf("Starting message listener for %s", clientID)
	}
	return nil, client
}
