package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store"
	"google.golang.org/protobuf/proto"
	"whatsgoingon/conf"
	"whatsgoingon/events"
	"whatsgoingon/routes"
)

func main() {
	conf.InitToken()
	r := gin.Default()
	store.DeviceProps.Os = proto.String("UatzAPI")

	err := godotenv.Load(".env")
	if err != nil {
		panic(errors.New("cannot open .env file"))
	}

	// Channel for init a listener in a goroutine
	go events.InitListener()

	r.Use(conf.CORSmiddleware())
	r.Use(conf.TokenMiddleware())

	// Device Routes
	r.GET("/connect", routes.DeviceNew)
	r.GET("/device", routes.DeviceList)
	r.GET("/device/:deviceId", routes.GetDeviceInfo)

	// Listener Routes
	r.GET("/start_listener", routes.StartListener)

	// Message Routes
	r.POST("/send/message", routes.SendMessage)
	r.POST("/send/sticker", routes.SendSticker)

	// Webhook Routes
	r.GET("/webhook", routes.WebhookList)
	r.POST("/webhook", routes.WebhookAdd)
	r.DELETE("/webhook/:deviceID", routes.WebhookRemove)
	r.GET("/webhook/:deviceID", routes.WebhookByDevice)
	r.GET("/webhook/:deviceID/all", routes.WebhookListByDevice)

	r.Run(":8080")
}
