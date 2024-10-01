package main

import (
	"log"

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
	// Set the device properties for WhatsApp API usage.
	store.DeviceProps.Os = proto.String("UatzAPI")

	// Load environment variables from .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, running in production mode")
		gin.SetMode(gin.ReleaseMode) // Set Gin mode to Release if no .env is found
	}

	// Initialize tokens for authentication or any necessary configuration.
	conf.InitToken()

	// Initialize the Gin router.
	r := gin.Default()

	// Start the event listener in a separate goroutine.
	go events.InitListener()

	// Add CORS and Token middlewares for handling requests.
	r.Use(conf.CORSmiddleware())
	r.Use(conf.TokenMiddleware())

	// Device Routes
	r.GET("/connect", routes.DeviceNew)              // Connect to a new device
	r.GET("/device", routes.DeviceList)              // Get a list of devices
	r.GET("/device/:deviceId", routes.GetDeviceInfo) // Get device information by device ID

	// Listener Routes
	r.GET("/start_listener", routes.StartListener) // Start listener for messages

	// Message Routes
	r.POST("/send/message", routes.SendMessage) // Send a text message
	r.POST("/send/sticker", routes.SendSticker) // Send a sticker

	// Webhook Routes
	r.GET("/webhook", routes.WebhookList)                       // List all webhooks
	r.POST("/webhook", routes.WebhookAdd)                       // Add a new webhook
	r.DELETE("/webhook/:deviceID", routes.WebhookRemove)        // Remove a webhook by device ID
	r.GET("/webhook/:deviceID", routes.WebhookByDevice)         // Get active webhook for a specific device
	r.GET("/webhook/:deviceID/all", routes.WebhookListByDevice) // List all webhooks for a device

	// Run the server on port 8080.
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
