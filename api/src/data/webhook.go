package data

import (
	"github.com/uptrace/bun"
	"time"
)

type WebhookMessage struct {
	bun.BaseModel `bun:"table:webhook_message,alias:wh_msg"`
	ID            int       `json:"id" bun:"id,pk,autoincrement"`
	Message       string    `json:"message" bun:"message,notnull"`
	Response      string    `json:"response" bun:"response,notnull"`
	CodeResponse  int       `json:"code_response" bun:"code_response,notnull"`
	Timestamp     time.Time `json:"timestamp" bun:"timestamp,notnull"`
}
