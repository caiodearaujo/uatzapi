package helpers

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ztrue/tracerr"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func SendMessage(jid string, message string, recipient string) (*whatsmeow.SendResponse, error) {
	client, err := GetWhatsAppClientByJID(jid)
	if err != nil {
		return nil, err
	}

	destination := types.JID{
		User:   recipient,
		Server: types.DefaultUserServer,
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