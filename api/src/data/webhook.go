package data

import (
	"time"

	"github.com/uptrace/bun"
)

type WebhookMessage struct {
	bun.BaseModel `bun:"table:webhook_message,alias:wh_msg"`
	ID            int       `json:"id" bun:"id,pk,autoincrement"`
	DeviceID	  string    `json:"device_id" bun:"device_id,notnull"`
	Message       string    `json:"message" bun:"message,notnull"`
	Response      string    `json:"response" bun:"response,notnull"`
	WebhookURL    string    `json:"webhook_url" bun:"webhook_url,notnull"`
	CodeResponse  int       `json:"code_response" bun:"code_response,notnull"`
	Timestamp     time.Time `json:"timestamp" bun:"timestamp,notnull"`
}
