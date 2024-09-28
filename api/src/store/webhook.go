package store

import (
	"context"
	"fmt"
	"time"
	"whatsgoingon/data"
	"whatsgoingon/handler"
)

const (
	NO_ROWS_IN_RESULT_SET = "sql: no rows in result set"
)

func GetWebhookURLs() ([]data.DeviceWebhook, error) {
	db := GetBunConnection()

	webhooks := []data.DeviceWebhook{}
	err := db.NewSelect().
		Model(&webhooks).
		Where("active = ?", true).
		Relation("Device").
		Scan(context.Background())

	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

func InactivateWebhookByDeviceID(deviceId int) (error, bool) {
	db := GetBunConnection()

	_, err := db.NewUpdate().
		Model(&data.DeviceWebhook{}).
		Set("active = ?", false).
		Where("device_id = ?", deviceId).
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("error inactivating webhook by device ID: %v", err), false
	}
	return nil, true
}

func CreateNewWebhook(deviceId int, webhookURL string) (error, bool) {
	device, err := GetDeviceById(deviceId)
	if err != nil {
		return fmt.Errorf("error getting device by ID: %v", err), false
	}

	err, _ = InactivateWebhookByDeviceID(deviceId)
	if err != nil {
		return fmt.Errorf("error inactivating webhook by device ID: %v", err), false
	}

	_, err = InsertIntoTable(&data.DeviceWebhook{
		DeviceID:   deviceId,
		Device:     &device,
		WebhookURL: webhookURL,
		Active:     true,
		Timestamp:  time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error inserting webhook into table: %v", err), false
	}

	return nil, true
}

func GetWebhookActiveByDeviceID(deviceId int) (data.DeviceWebhook, error) {
	db := GetBunConnection()

	webhook := data.DeviceWebhook{}
	err := db.NewSelect().
		Model(&webhook).
		Where("device_id = ? AND active = ?", deviceId, true).
		Scan(context.Background())

	if err != nil {
		return webhook, err
	}
	return webhook, nil
}

func GetWebhooksByDeviceID(deviceId int) ([]data.DeviceWebhook, error) {
	db := GetBunConnection()

	webhooks := []data.DeviceWebhook{}
	err := db.NewSelect().
		Model(&webhooks).
		Where("device_id = ?", deviceId).
		Scan(context.Background())

	if err != nil {
		return webhooks, err
	}
	return webhooks, nil
}

// GetWebhookURLFordeviceID retrieves the webhook URL for the given device ID.
func GetWebhookURLByDeviceID(deviceID int) (string, bool, error) {
	db := GetBunConnection()

	deviceWebhook := new(data.DeviceWebhook)
	err := db.NewSelect().
		Model(deviceWebhook).
		Where("id = ? AND active = ?", deviceID, true).
		Scan(context.Background())

	if err != nil {
		handler.FailOnError(err, "Failed to get webhook URL for device ID")
		return "", false, err
	}
	return deviceWebhook.WebhookURL, deviceWebhook.Active, nil
}
