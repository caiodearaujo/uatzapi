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

func SendMessage(jid string, message string, recipient string) (*whatsmeow.SendResponse, error) {
	client, err := GetWhatsAppClientByJID(jid)
	if err != nil {
		return nil, err
	}

	destination, err := CheckIfNumberExistsAndGetJID(recipient, client)
	if err != nil {
		return nil, err
	}

	encryptedMessage := &waE2E.Message{
		Conversation: proto.String(message),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Timeout context
	defer cancel()

	resp, err := client.SendMessage(ctx, destination, encryptedMessage)
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrMessageSending, err))
	}

	return &resp, nil
}

func SendSticker(jid string, stickerData []byte, recipient string) (*whatsmeow.SendResponse, error) {
	client, err := GetWhatsAppClientByJID(jid)
	if err != nil {
		return nil, err
	}

	destination, err := CheckIfNumberExistsAndGetJID(recipient, client)
	if err != nil {
		return nil, err
	}

	stickerDataWebp, err := ConvertImageToWebp(stickerData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert sticker to webp: %v", err)
	}

	stickerUpload, err := client.Upload(context.Background(), stickerDataWebp, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("failed to upload sticker: %v", err)
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Timeout context
	defer cancel()

	resp, err := client.SendMessage(ctx, destination, encryptedMessage)
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrMessageSending, err))
	}

	return &resp, nil
}
