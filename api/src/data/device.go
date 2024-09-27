package data

import (
	"strconv"
	"time"

	"github.com/uptrace/bun"
)

type Device struct {
	bun.BaseModel 		   `bun:"table:device,alias:dvc"`
	ID           int       `json:"id" bun:"id,pk,autoincrement"`
	JID     	 string    `json:"whatsapp_id" bun:"whatsapp_id,notnull,unique"`
	PushName     string    `json:"push_name" bun:"push_name,notnull"`
	BusinessName string    `json:"business_name" bun:"business_name,"`
	Active	   	 bool      `json:"active" bun:"active,notnull"`
	CreatedAt   time.Time `json:"created_at" bun:"created_at,notnull"`
}

type DeviceHandler struct {
	bun.BaseModel 		  `bun:"table:device_handler,alias:dvc_hdl"`
	ID          int       `json:"id" bun:"id,pk,autoincrement"`
	DeviceID    string    `json:"device_id" bun:"device_id,notnull"`
	Active 		bool      `json:"active" bun:"active,notnull"`
	ActiveAt    time.Time `json:"timestamp" bun:"timestamp,notnull"`
	InactiveAt  time.Time `json:"inactive_at" bun:"inactive_at"`
	Device	    *Device    `bun:"rel:belongs-to,join:device_id=id"`
}

type DeviceWebhook struct {
	bun.BaseModel           `bun:"table:device_webhook,alias:dvc_wh"`
	ID		      int       `json:"id" bun:"id,pk,autoincrement"`	
	DeviceID      string    `json:"device_id" bun:"device_id,notnull"`
	WebhookURL    string    `json:"webhook_url" bun:"webhook_url,notnull"`
	Active        bool      `json:"active" bun:"active,notnull"`
	Timestamp     time.Time `json:"timestamp" bun:"timestamp,notnull"`
	Device	      *Device   `bun:"rel:belongs-to,join:device_id=id"`
}

func (m Device) DeviceID() string {
	return strconv.Itoa(m.ID)
}
