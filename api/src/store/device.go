package store

import (
	"context"
	"fmt"
	"time"
	"whatsgoingon/data"
)

// GetJIDByDeviceID retrieves the JID (WhatsApp ID) for a given device ID.
func GetJIDByDeviceID(deviceID int) (string, error) {
	db := GetBunConnection()

	device := new(data.Device)
	err := db.NewSelect().
		Model(device).
		Where("id = ? AND active = ?", deviceID, true).
		Scan(context.Background())

	if err != nil {
		return "", fmt.Errorf("failed to retrieve JID for device ID %d: %v", deviceID, err)
	}
	return device.JID, nil
}

// GetDeviceByID retrieves the device details by its ID.
func GetDeviceByID(deviceID int) (data.Device, error) {
	db := GetBunConnection()

	device := new(data.Device)
	err := db.NewSelect().
		Model(device).
		Where("id = ? AND active = ?", deviceID, true).
		Scan(context.Background())

	if err != nil {
		return data.Device{}, fmt.Errorf("failed to retrieve device by ID %d: %v", deviceID, err)
	}
	return *device, nil
}

// GetDeviceByJID retrieves the device details by its WhatsApp ID (JID).
func GetDeviceByJID(jid string) (data.Device, error) {
	db := GetBunConnection()

	device := new(data.Device)
	err := db.NewSelect().
		Model(device).
		Where("whatsapp_id = ? AND active = ?", jid, true).
		Scan(context.Background())

	if err != nil {
		return data.Device{}, fmt.Errorf("failed to retrieve device by JID %s: %v", jid, err)
	}
	return *device, nil
}

// InsertDeviceIfNotExists inserts a new device into the database if it doesn't already exist.
func InsertDeviceIfNotExists(device *data.Device) (*data.Device, error) {
	db := GetBunConnection()

	// Check if the device already exists by JID
	exists, err := db.NewSelect().
		Model(device).
		Where("whatsapp_id = ?", device.JID).
		Exists(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to check if device exists: %v", err)
	}

	if !exists {
		// Insert the new device into the table
		_, err := db.NewInsert().Model(device).Returning("*").Exec(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to insert device into table: %v", err)
		}
		fmt.Printf("Device inserted successfully with ID: %d\n", device.ID)
		return device, nil
	}
	return nil, fmt.Errorf("device with JID %s already exists", device.JID)
}

// BulkUpdateDeviceHandlerOff deactivates all active device handlers in the database.
func BulkUpdateDeviceHandlerOff() error {
	db := GetBunConnection()

	// Set all active handlers to inactive
	_, err := db.NewUpdate().
		Model((*data.DeviceHandler)(nil)).
		Set("active = false").
		Set("inactive_at = ?", time.Now()).
		Where("active = true").
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("failed to deactivate device handlers: %v", err)
	}

	fmt.Println("All active device handlers have been deactivated.")
	return nil
}

// GetTop20WebhookMessagesByDeviceID retrieves the last 20 webhook messages for a given device.
func GetTop20WebhookMessagesByDeviceID(deviceID int) ([]data.WebhookMessage, error) {
	db := GetBunConnection()

	var webhookMessages []data.WebhookMessage
	err := db.NewSelect().
		Model(&webhookMessages).
		Where("device_id = ?", deviceID).
		Order("timestamp DESC").
		Limit(20).
		Scan(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve top 20 webhook messages for device ID %d: %v", deviceID, err)
	}

	return webhookMessages, nil
}

// InactiveWebhookURLByDeviceID deactivates the webhook URL for a specific device ID.
func InactiveWebhookURLByDeviceID(deviceID int) error {
	db := GetBunConnection()

	_, err := db.NewUpdate().
		Model((*data.DeviceWebhook)(nil)).
		Set("active = false").
		Where("device_id = ?", deviceID).
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("failed to deactivate webhook URL for device ID %d: %v", deviceID, err)
	}

	fmt.Printf("Webhook URL for device ID %d has been deactivated.\n", deviceID)
	return nil
}

// RemoveDevice removes a device from the database by its ID.
func RemoveDevice(deviceID int) error {
	db := GetBunConnection()

	// Fetch the device by ID to ensure it exists
	device := new(data.Device)
	err := db.NewSelect().
		Model(device).
		Where("id = ?", deviceID).
		Scan(context.Background())
	if err != nil {
		return fmt.Errorf("failed to find device with ID %d: %v", deviceID, err)
	}

	// Delete the device from the database
	_, err = db.NewDelete().
		Model(device).
		Where("id = ?", deviceID).
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete device with ID %d: %v", deviceID, err)
	}

	fmt.Printf("Device with ID %d has been successfully removed.\n", deviceID)
	return nil
}
