package store

import (
	"context"
	"fmt"
	"time"
	"whatsgoingon/data"
	"whatsgoingon/handler"
)

// GetWebhookURLs retrieves all active webhook URLs from the database.
func GetWebhookURLs() ([]data.DeviceWebhook, error) {
	db := GetBunConnection()

	var webhooks []data.DeviceWebhook
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

// InactivateWebhookByDeviceID sets the 'active' field of a webhook to false based on the device ID.
// It returns an error if the operation fails, along with a boolean indicating success.
func InactivateWebhookByDeviceID(deviceID int) (error, bool) {
	db := GetBunConnection()

	_, err := db.NewUpdate().
		Model(&data.DeviceWebhook{}).
		Set("active = ?", false).
		Where("device_id = ?", deviceID).
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("error inactivating webhook for device ID %d: %v", deviceID, err), false
	}
	return nil, true
}

// CreateNewWebhook creates a new webhook for the given device ID.
// It inactivates any existing webhooks for the device before inserting a new one.
func CreateNewWebhook(deviceID int, webhookURL string) (error, bool) {
	device, err := GetDeviceById(deviceID)
	if err != nil {
		return fmt.Errorf("error retrieving device by ID %d: %v", deviceID, err), false
	}

	err, _ = InactivateWebhookByDeviceID(deviceID)
	if err != nil {
		return fmt.Errorf("error inactivating webhook for device ID %d: %v", deviceID, err), false
	}

	_, err = InsertIntoTable(&data.DeviceWebhook{
		DeviceID:   deviceID,
		Device:     &device,
		WebhookURL: webhookURL,
		Active:     true,
		Timestamp:  time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error inserting webhook for device ID %d: %v", deviceID, err), false
	}

	return nil, true
}

// GetWebhookActiveByDeviceID retrieves the active webhook for the given device ID.
func GetWebhookActiveByDeviceID(deviceID int) (data.DeviceWebhook, error) {
	db := GetBunConnection()

	var webhook data.DeviceWebhook
	err := db.NewSelect().
		Model(&webhook).
		Where("device_id = ? AND active = ?", deviceID, true).
		Scan(context.Background())

	if err != nil {
		return webhook, err
	}
	return webhook, nil
}

// GetWebhooksByDeviceID retrieves all webhooks associated with a given device ID.
func GetWebhooksByDeviceID(deviceID int) ([]data.DeviceWebhook, error) {
	db := GetBunConnection()

	var webhooks []data.DeviceWebhook
	err := db.NewSelect().
		Model(&webhooks).
		Where("device_id = ?", deviceID).
		Scan(context.Background())

	if err != nil {
		return webhooks, err
	}
	return webhooks, nil
}

// GetWebhookURLByDeviceID retrieves the webhook URL for the given device ID if it is active.
// Returns the URL, the active status, and an error if applicable.
func GetWebhookURLByDeviceID(deviceID int) (string, bool, error) {
	db := GetBunConnection()

	var deviceWebhook data.DeviceWebhook
	err := db.NewSelect().
		Model(&deviceWebhook).
		Where("id = ? AND active = ?", deviceID, true).
		Scan(context.Background())

	if err != nil {
		handler.FailOnError(err, "Failed to get webhook URL for device ID")
		return "", false, err
	}
	return deviceWebhook.WebhookURL, deviceWebhook.Active, nil
}
