package data

import (
	"github.com/uptrace/bun"
	"time"
)

type Device struct {
	ID           int       `json:"id" bun:"id,pk,autoincrement"`
	DeviceID     string    `json:"device_id" bun:"device_id,notnull"`
	PushName     string    `json:"push_name" bun:"push_name,notnull"`
	BusinessName string    `json:"business_name" bun:"business_name,"`
	Timestamp    time.Time `json:"timestamp" bun:"timestamp,notnull"`
}

type DeviceHandler struct {
	bun.BaseModel `bun:"table:device_handler,alias:dev_hdl"`
	Device
	Active bool `json:"active" bun:"active,notnull"`
}

type DeviceWebhook struct {
	bun.BaseModel `bun:"table:device_webhook,alias:dev_wh"`
	DeviceID      string    `json:"device_id" bun:"device_id,notnull"`
	WebhookURL    string    `json:"webhook_url" bun:"webhook_url,notnull"`
	Active        bool      `json:"active" bun:"active,notnull"`
	Timestamp     time.Time `json:"timestamp" bun:"timestamp,notnull"`
}
