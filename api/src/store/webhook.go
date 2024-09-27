package store

import (
	"context"
	"whatsgoingon/data"
)

func GetWebhookURLs() ([]data.DeviceWebhook, error) {
	db := GetBunConnection()

	webhooks := []data.DeviceWebhook{}
	err := db.NewSelect().
		Model(&webhooks).
		Relation("Device").
		Scan(context.Background())

	if err != nil {
		return nil, err
	}
	return webhooks, nil
}