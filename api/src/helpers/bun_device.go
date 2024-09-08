package helpers

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"whatsgoingon/data"
)

func InsertIntoTable(model interface{}) (interface{}, error) {
	pgDB := getPostgresConnection()
	db := bun.NewDB(pgDB, pgdialect.New())
	res, err := db.NewInsert().Model(model).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func BulkUpdateDeviceHandlerOff() {
	pgDB := getPostgresConnection()
	db := bun.NewDB(pgDB, pgdialect.New())

	_, err := db.NewUpdate().NewUpdate().Model((*data.DeviceHandler)(nil)).Set("active = false").Where("active = true").Exec(context.Background())
	if err != nil {
		failOnError(err, "Failed to update table")
	} else {
		fmt.Printf("Table updated successfully")
	}
}

func GetWebhookURLForClientID(clientID string) (string, error) {
	db := GetBunConnection()

	deviceWebhook := new(data.DeviceWebhook)
	err := db.NewSelect().Model(deviceWebhook).Where("device_id = ?", clientID).Scan(context.Background())
	if err != nil {
		return "", err
	} else {
		return deviceWebhook.WebhookURL, nil
	}
}
