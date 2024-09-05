package store

import (
	"fmt"
	"whatsgoingon/data"

	"github.com/ztrue/tracerr"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type Conversation struct {
	Side     string `json:"side"`
	Message  string `json:"message"`
	MimeType string `json:"mime_type"`
}

// SaveMessage stores a message in the database
func SaveMessage(msg events.Message, client *whatsmeow.Client) (error, *data.StoredMessage) {
	content, err := data.ConvertEventToStoredMessage(msg, client)
	if err != nil {
		tracerr.Print(err)
		return fmt.Errorf("error getting message content: %v", err), nil
	}

	// Save the message to the database
	

	return nil, &content
}
