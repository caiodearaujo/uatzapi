package data

import (
	"errors"
	"log"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// StoredMessage represents a message stored in the database.
type StoredMessage struct {
	MessageID       string    `firestore:"message_id" json:"message_id"`               // Unique identifier for the message
	IsFromMe        bool      `firestore:"is_from_me" json:"is_from_me"`               // Whether the message was sent by the current user
	IsFromGroup     bool      `firestore:"is_from_group" json:"is_from_group"`         // Whether the message is from a group chat
	MediaType       string    `firestore:"media_type" json:"media_type"`               // Type of message (TEXT, IMAGE, VIDEO, etc.)
	Text            string    `firestore:"text" json:"text"`                           // Text content of the message (if applicable)
	Content         []byte    `firestore:"content" json:"content"`                     // Raw content of the message (for media)
	ContentMimeType string    `firestore:"content_mime_type" json:"content_mime_type"` // Mime type of the content media
	RecipientID     string    `firestore:"recipient_id" json:"recipient_id"`           // WhatsApp ID of the recipient
	RecipientName   string    `firestore:"push_name" json:"push_name"`                 // Display name of the recipient
	Timestamp       time.Time `firestore:"timestamp" json:"timestamp"`                 // Timestamp of the message
}

// ConvertEventToStoredMessage converts a WhatsApp event message into a StoredMessage structure.
// It extracts media or text content and assigns the appropriate fields in the StoredMessage.
func ConvertEventToStoredMessage(v events.Message, client *whatsmeow.Client) (*StoredMessage, error) {
	messageContent := StoredMessage{
		MessageID:     v.Info.ID,        // Unique identifier for the message
		RecipientID:   v.Info.Chat.User, // WhatsApp ID of the recipient
		Timestamp:     v.Info.Timestamp, // Timestamp when the message was sent/received
		IsFromMe:      v.Info.IsFromMe,  // Indicates if the message was sent by the user
		IsFromGroup:   v.Info.IsGroup,   // Indicates if the message belongs to a group chat
		RecipientName: v.Info.PushName,  // Display name of the recipient
	}

	// Handle Image Messages
	if v.Message.ImageMessage != nil {
		messageContent.MediaType = "IMAGE"
		if v.Message.ImageMessage.Caption != nil {
			messageContent.Text = v.Message.GetImageMessage().GetCaption()
		}
		messageContent.ContentMimeType = v.Message.ImageMessage.GetMimetype()
		content, err := client.Download(v.Message.ImageMessage)
		messageContent.Content = content
		return &messageContent, err
	}

	// Handle Video Messages
	if v.Message.VideoMessage != nil {
		messageContent.MediaType = "VIDEO"
		messageContent.ContentMimeType = v.Message.VideoMessage.GetMimetype()
		content, err := client.Download(v.Message.VideoMessage)
		messageContent.Content = content
		return &messageContent, err
	}

	// Handle Audio Messages
	if v.Message.AudioMessage != nil {
		messageContent.MediaType = "AUDIO"
		messageContent.ContentMimeType = v.Message.AudioMessage.GetMimetype()
		content, err := client.Download(v.Message.AudioMessage)
		messageContent.Content = content
		messageContent.Text = "" // Placeholder for potential speech-to-text functionality
		return &messageContent, err
	}

	// Handle Sticker Messages
	if v.Message.StickerMessage != nil {
		messageContent.MediaType = "STICKER"
		messageContent.ContentMimeType = v.Message.StickerMessage.GetMimetype()
		content, err := client.Download(v.Message.StickerMessage)
		messageContent.Content = content
		return &messageContent, err
	}

	// Handle Document Messages
	if v.Message.DocumentMessage != nil {
		messageContent.MediaType = "DOCUMENT"
		messageContent.ContentMimeType = v.Message.DocumentMessage.GetMimetype()
		content, err := client.Download(v.Message.DocumentMessage)
		messageContent.Content = content
		return &messageContent, err
	}

	// Handle Text and Extended Text Messages
	if v.Message.Conversation != nil {
		messageContent.MediaType = "TEXT"
		messageContent.Text = v.Message.GetConversation()
		messageContent.ContentMimeType = "text/plain"
		return &messageContent, nil
	}

	if v.Message.ExtendedTextMessage != nil {
		messageContent.MediaType = "TEXT"
		messageContent.Text = v.Message.GetExtendedTextMessage().GetText()
		messageContent.ContentMimeType = "text/plain"
		return &messageContent, nil
	}

	// Handle Unrecognized Message Types
	if v.Message != nil {
		messageContent.MediaType = "UNKNOWN"
		messageContent.ContentMimeType = "text/plain"
		messageContent.Text = v.Message.String()
	} else {
		log.Printf("Unrecognized message type: %T", v.Message)
		return &messageContent, errors.New("message type not recognized")
	}

	return &messageContent, nil
}
