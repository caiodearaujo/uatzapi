package main

import (
	"errors"
	"whatsgoingon/events"
	"whatsgoingon/routes"
	myStore "whatsgoingon/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store"
	"google.golang.org/protobuf/proto"
)

func main() {
	r := gin.Default()
	store.DeviceProps.Os = proto.String("UatzAPI")

	err := godotenv.Load(".env")
	if err != nil {
		panic(errors.New("cannot open .env file"))
	}

	// Channel for init a listener in a goroutine
	myStore.CreateTablesFromDataPkg()
	go events.InitListener()

	r.Use(cors.Default())

	r.GET("/connect", routes.Connect)

	r.GET("/devices", routes.Devices)

	r.GET("/start_listener", routes.StartListener)

	r.GET("/send-message", routes.SendMessage)

	r.GET("/webhook", routes.WebhookList)

	r.POST("/webhook", routes.WebhookAdd)

	r.DELETE("/webhook/:deviceID", routes.WebhookRemove)

	r.GET("/historical", routes.HistoricalConversation)

	r.Run(":8080")
}
