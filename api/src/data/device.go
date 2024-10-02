package data

import (
	"time"

	"github.com/uptrace/bun"
)

// Device represents a WhatsApp device that is linked to the system.
type Device struct {
	bun.BaseModel `bun:"table:device,alias:dvc"` // Specifies the table name and alias
	ID            int                            `json:"id" bun:"id,pk,autoincrement"`                 // Primary key, auto-incremented
	JID           string                         `json:"whatsapp_id" bun:"whatsapp_id,notnull,unique"` // WhatsApp JID (WhatsApp ID), must be unique
	PushName      string                         `json:"push_name" bun:"push_name,notnull"`            // Display name of the WhatsApp device
	BusinessName  string                         `json:"business_name" bun:"business_name"`            // Optional business name for WhatsApp Business API
	Active        bool                           `json:"active" bun:"active,notnull"`                  // Indicates if the device is currently active
	CreatedAt     time.Time                      `json:"created_at" bun:"created_at,notnull"`          // Timestamp when the device was added
}

// DeviceHandler represents a record of when a device's handler is active or inactive.
type DeviceHandler struct {
	bun.BaseModel `bun:"table:device_handler,alias:dvc_hdl"` // Specifies the table name and alias
	ID            int                                        `json:"id" bun:"id,pk,autoincrement"`      // Primary key, auto-incremented
	DeviceID      int                                        `json:"device_id" bun:"device_id,notnull"` // Foreign key to the device table
	Active        bool                                       `json:"active" bun:"active,notnull"`       // Indicates if the handler is currently active
	ActiveAt      time.Time                                  `json:"active_at" bun:"active_at,notnull"` // Timestamp when the handler became active
	InactiveAt    time.Time                                  `json:"inactive_at" bun:"inactive_at"`     // Timestamp when the handler became inactive (nullable)
	Device        *Device                                    `bun:"rel:belongs-to,join:device_id=id"`   // Relation to the device table (many-to-one)
}

// DeviceWebhook represents a webhook URL associated with a device for receiving updates.
type DeviceWebhook struct {
	bun.BaseModel `bun:"table:device_webhook,alias:dvc_wh"` // Specifies the table name and alias
	ID            int                                       `json:"id" bun:"id,pk,autoincrement"`          // Primary key, auto-incremented
	DeviceID      int                                       `json:"device_id" bun:"device_id,notnull"`     // Foreign key to the device table
	WebhookURL    string                                    `json:"webhook_url" bun:"webhook_url,notnull"` // The webhook URL to send updates to
	Active        bool                                      `json:"active" bun:"active,notnull"`           // Indicates if the webhook is currently active
	Timestamp     time.Time                                 `json:"timestamp" bun:"timestamp,notnull"`     // Timestamp when the webhook was created
	Device        *Device                                   `bun:"rel:belongs-to,join:device_id=id"`       // Relation to the device table (many-to-one)
}
