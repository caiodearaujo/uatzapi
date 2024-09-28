package store

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"whatsgoingon/data"
	"whatsgoingon/handler"
)

// GetDeviceIDBydeviceID retrieves the device ID by the device ID.
func GetJIDByDeviceID(deviceID int) (string, error) {
	db := GetBunConnection()

	device := new(data.Device)
	err := db.NewSelect().
		Model(device).
		Where("id = ? AND active = ?", deviceID, true).
		Scan(context.Background())

	if err != nil {
		return "", err
	}
	return device.JID, nil
}

func GetDeviceById(deviceID int) (data.Device, error) {
	db := GetBunConnection()

	device := new(data.Device)
	err := db.NewSelect().
		Model(device).
		Where("id = ? AND active = ?", deviceID, true).
		Scan(context.Background())

	if err != nil {
		return data.Device{}, err
	}
	return *device, nil
}

func GetDeviceByJID(jid string) (data.Device, error) {
	db := GetBunConnection()

	device := new(data.Device)
	err := db.NewSelect().
		Model(device).
		Where("whatsapp_id = ? AND active = ?", jid, true).
		Scan(context.Background())

	if err != nil {
		return data.Device{}, err
	}
	return *device, nil
}

// InsertIntoTableIfNotExists inserts a given model into the database if it does not exist.
func InsertDeviceIfNotExists(device *data.Device) (*data.Device, error) {
	db := GetBunConnection()

	// Verify is exists a device with the same ID.
	exists, err := db.NewSelect().
		Model(device).
		Where("whatsapp_id = ?", device.JID).
		Exists(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to verify if device exists: %v", err)
	}

	if !exists {
		// Insert the model into the table.
		_, err := db.NewInsert().Model(device).Returning("*").Exec(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to insert device into table: %v", err)
		}
		fmt.Printf("Device inserted successfully " + strconv.FormatInt(int64(device.ID), 10))
		return device, nil
	}
	return nil, fmt.Errorf("device already exists")
}

// BulkUpdateDeviceHandlerOff deactivates all active device handlers.
func BulkUpdateDeviceHandlerOff() error {
	db := GetBunConnection()

	// Update all active device handlers do inactive.
	_, err := db.NewUpdate().
		Model((*data.DeviceHandler)(nil)).
		Set("active = false").
		Set("inactived_at = ?", time.Now()).
		Where("active = true").
		Exec(context.Background())

	if err != nil {
		handler.FailOnError(err, "Failed to update table")
		return err
	}

	fmt.Printf("Device handlers bulk updated (inactive) successfully")
	return nil
}

// GetTop20WebhookMessages retrieves the top 20 webhook messages.
func GetTop20WebhookMessagesByDeviceID(deviceID int) []data.WebhookMessage {
	db := GetBunConnection()

	var webhookMessages []data.WebhookMessage
	err := db.NewSelect().
		Model(&webhookMessages).
		Where("id = ?", deviceID).
		Order("timestamp DESC").
		Limit(20).
		Scan(context.Background())

	if err != nil {
		handler.FailOnError(err, "Failed to get top 20 webhook messages")
		return []data.WebhookMessage{}
	}

	return webhookMessages
}

// InactiveWebhookURLFordeviceID deactivates the webhook URL for the given device ID.
func InactiveWebhookURLByDeviceID(deviceID int) error {
	db := GetBunConnection()

	_, err := db.NewUpdate().
		Model((*data.DeviceWebhook)(nil)).
		Set("active = false").
		Where("id = ?", deviceID).
		Exec(context.Background())

	if err != nil {
		handler.FailOnError(err, "Failed to inactivate webhook URL for device ID")
		return err
	}

	fmt.Printf("Webhook URL for device ID %s inactivated successfully", deviceID)
	return nil
}

func RemoveDevice(deviceID int) (error, bool) {
	db := GetBunConnection()

	device := new(data.Device)
	err := db.NewSelect().
		Model(device).
		Where("id = ?", deviceID).
		Scan(context.Background())

	if err != nil {
		return err, false
	}

	_, err = db.NewDelete().
		Model(device).
		Where("id = ?", deviceID).
		Exec(context.Background())

	if err != nil {
		return err, false
	}

	return nil, true
}
