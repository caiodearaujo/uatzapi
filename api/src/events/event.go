package events

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ztrue/tracerr"
	"go.mau.fi/whatsmeow/types/events"
	"log"
	"whatsgoingon/data"
	"whatsgoingon/helpers"
	myStore "whatsgoingon/store"
)

func InitListener() {
	clientIds, err := helpers.GetAllClientIDs()
	failOnError(err, "Get clientIds failed")

	for _, clientId := range clientIds {
		StartMessageListener(clientId)
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func StartMessageListener(clientID string) {
	client, err := helpers.GetClientById(clientID)
	if err != nil {
		tracerr.Print(err)
	} else {
		client.AddEventHandler(func(evt interface{}) {
			log.Printf("Event received")
			switch v := evt.(type) {
			case *events.Message:
				if err, content := myStore.SaveMessage(*v, client); err != nil {
					tracerr.Print(err)
				} else {
					sendToRabbitMQ(*content, clientID) // async
				}
			}
		})
		log.Printf("Starting message listener for %s", clientID)
	}
}

func sendToRabbitMQ(content data.StoredMessage, clientID string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		clientID,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	body, err := json.Marshal(content)
	failOnError(err, "Failed to encode a message")

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent to client %s", clientID)
}
