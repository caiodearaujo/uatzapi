package helpers

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ztrue/tracerr"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

// SendMessage sends a text message to a recipient using WhatsApp.
// The function retrieves the WhatsApp client by JID, checks if the recipient number exists, and sends the message.
func SendMessage(jid string, message string, recipient string) (*whatsmeow.SendResponse, error) {
	// Retrieve WhatsApp client for the given JID.
	client, err := GetWhatsAppClientByJID(jid)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect()

	// Check if the recipient number exists and get its JID.
	destination, err := CheckIfNumberExistsAndGetJID(recipient, client)
	if err != nil {
		return nil, err
	}

	// Create an encrypted WhatsApp message.
	encryptedMessage := &waE2E.Message{
		Conversation: proto.String(message),
	}

	// Set a timeout for the context to avoid blocking.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send the message and return the response or error.
	resp, err := client.SendMessage(ctx, destination, encryptedMessage)
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrMessageSending, err))
	}

	return &resp, nil
}

// SendSticker sends a sticker to a recipient on WhatsApp.
// It retrieves the WhatsApp client, converts the sticker image to WebP, uploads it, and sends the sticker message.
func SendSticker(jid string, stickerData []byte, recipient string) (*whatsmeow.SendResponse, error) {
	// Retrieve WhatsApp client for the given JID.
	client, err := GetWhatsAppClientByJID(jid)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect()

	// Check if the recipient number exists and get its JID.
	destination, err := CheckIfNumberExistsAndGetJID(recipient, client)
	if err != nil {
		return nil, err
	}

	// Convert the sticker image to WebP format.
	stickerDataWebp, err := ConvertImageToWebp(stickerData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert sticker to webp: %v", err)
	}

	// Upload the WebP sticker to WhatsApp servers.
	stickerUpload, err := client.Upload(context.Background(), stickerDataWebp, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("failed to upload sticker: %v", err)
	}

	// Create the sticker message with the upload details.
	encryptedMessage := &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           proto.String(stickerUpload.URL),
			DirectPath:    proto.String(stickerUpload.DirectPath),
			MediaKey:      stickerUpload.MediaKey,
			Mimetype:      proto.String("image/webp"),
			FileSHA256:    stickerUpload.FileSHA256,
			FileEncSHA256: stickerUpload.FileEncSHA256,
			FileLength:    proto.Uint64(stickerUpload.FileLength),
		},
	}

	// Set a timeout for the context to avoid blocking.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send the sticker message and return the response or error.
	resp, err := client.SendMessage(ctx, destination, encryptedMessage)
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrMessageSending, err))
	}

	return &resp, nil
}
