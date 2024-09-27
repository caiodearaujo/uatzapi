package routes

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"net/http"
	"whatsgoingon/helpers"
	"whatsgoingon/events"
)

func Connect(c *gin.Context) {
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
			qrCodeBase64, err := generateQRCode(evt.Code)
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

func generateQRCode(qrCode string) (string, error) {
	// Gera a imagem do QR code
	png, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to encode QR code: %w", err)
	}

	// Converte o PNG para base64
	buf := new(bytes.Buffer)
	if _, err := buf.Write(png); err != nil { // Corrigido aqui
		return "", fmt.Errorf("failed to write PNG to buffer: %w", err)
	}
	qrCodeBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return qrCodeBase64, nil
}
