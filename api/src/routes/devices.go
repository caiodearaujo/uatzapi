package routes

import (
	"context"
	"net/http"
	"whatsgoingon/events"
	"whatsgoingon/helpers"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func DeviceList(c *gin.Context) {
	devices, _ := helpers.GetDeviceList()

	c.JSON(http.StatusOK, devices)
}

func DeviceNew(c *gin.Context) {
	var err error
	client, err := helpers.NewClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client.Disconnect()
	qrChan, _ := client.GetQRChannel(context.Background())
	err = client.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for evt := range qrChan {
		if evt.Event == "code" {
			qrCodeBase64, err := helpers.GenerateQRCode(evt.Code)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			client.AddEventHandler(events.NewClientHandler(client))
			c.JSON(http.StatusOK, gin.H{"qrCode": qrCodeBase64})
			return
		}
	}

}
