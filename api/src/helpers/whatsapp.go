package helpers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"whatsgoingon/handler"
	"whatsgoingon/store"

	"github.com/skip2/go-qrcode"
	"github.com/ztrue/tracerr"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Error definitions for database and WhatsApp client issues.
var (
	dbLog                 = waLog.Stdout("Database", "WARN", true)
	wmLog                 = waLog.Stdout("WhatsMeow", "WARN", true)
	dbMutex               sync.Mutex
	ErrDBConnectionFailed = errors.New("failed to connect to the database")
	ErrDeviceNotFound     = errors.New("device not found in the store")
	ErrClientConnection   = errors.New("failed to connect the WhatsApp client")
	ErrMessageSending     = errors.New("failed to send the message")
)

// DeviceResponse represents the device information returned in API responses.
type DeviceResponse struct {
	ID           int       `json:"id"`
	Number       string    `json:"number"`
	PushName     string    `json:"push_name"`
	BusinessName string    `json:"business_name"`
	Contacts     int       `json:"contacts"`
	Timestamp    time.Time `json:"timestamp"`
}

// ContactsResponse represents the contact information associated with a device.
type ContactsResponse struct {
	Name           string `json:"name"`
	Number         string `json:"number"`
	ProfilePicture string `json:"profile_picture"`
}

// DeviceInfoResponse is the detailed response structure for device information.
type DeviceInfoResponse struct {
	DeviceID         int                `json:"device_id"`
	ProfilePicture   string             `json:"profile_picture"`
	Webhook          string             `json:"webhook"`
	PhoneNumber      string             `json:"phone_number"`
	PushName         string             `json:"push_name"`
	BusinessName     string             `json:"business_name"`
	ContactsResponse []ContactsResponse `json:"contacts"`
}

// connectToDatabase establishes a connection to the PostgreSQL database and returns a WhatsMeow container.
func connectToDatabase() (*sqlstore.Container, error) {
	dbUser := os.Getenv("PG_USERNAME")
	dbPwd := os.Getenv("PG_PASSWORD")
	dbTCPHost := os.Getenv("PG_HOSTNAME")
	dbPort := os.Getenv("PG_PORT")
	dbName := os.Getenv("PG_DATABASE")
	dbSchema := os.Getenv("PG_WM_SCHEMA")

	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s search_path=%s sslmode=disable",
		dbTCPHost, dbUser, dbPwd, dbPort, dbName, dbSchema)

	container, err := sqlstore.New("pgx", dbURI, dbLog)
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrDBConnectionFailed, err))
	}
	return container, nil
}

// GetWhatsAppClientByJID retrieves a WhatsApp client by its JID (WhatsApp ID).
func GetWhatsAppClientByJID(whatsappID string) (*whatsmeow.Client, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	container, err := connectToDatabase()
	if err != nil {
		return nil, err
	}

	jid, _ := types.ParseJID(whatsappID)

	deviceStore, err := container.GetDevice(jid)
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrDeviceNotFound, err))
	}

	client := whatsmeow.NewClient(deviceStore, wmLog)
	if client.Store.ID == nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrClientConnection, client.Store.ID))
	}

	if !client.IsConnected() {
		err := client.Connect()
		if err != nil {
			return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrClientConnection, err))
		}
	}
	return client, nil
}

// GetWhatsappClientByDeviceID returns a WhatsApp client using the device's ID.
func GetWhatsappClientByDeviceID(deviceID int) (*whatsmeow.Client, error) {
	device, err := store.GetDeviceByID(deviceID)
	if err != nil {
		return nil, err
	}
	return GetWhatsAppClientByJID(device.JID)
}

// GetAllWhatsappIDs retrieves all WhatsApp IDs from the database.
func GetAllWhatsappIDs() ([]string, error) {
	container, err := connectToDatabase()
	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetAllDevices()
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrDeviceNotFound, err))
	}

	var deviceIDs []string
	for _, device := range deviceStore {
		deviceIDs = append(deviceIDs, device.ID.String())
	}

	return deviceIDs, nil
}

// NewClient creates a new WhatsApp client instance.
func NewClient() (*whatsmeow.Client, error) {
	container, err := connectToDatabase()
	if err != nil {
		return nil, err
	}

	deviceStore := container.NewDevice()
	client := whatsmeow.NewClient(deviceStore, wmLog)

	return client, nil
}

// GetDeviceList retrieves a list of all registered devices and their associated information.
func GetDeviceList() ([]DeviceResponse, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	container, err := connectToDatabase()
	if err != nil {
		return nil, err
	}

	devices, err := container.GetAllDevices()
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrDeviceNotFound, err))
	}

	var deviceList []DeviceResponse
	for _, device := range devices {
		allContacts, err := device.Contacts.GetAllContacts()
		if err != nil {
			continue // Log error and skip to the next device
		}

		if dvc, err := store.GetDeviceByJID(device.ID.String()); err == nil {
			deviceList = append(deviceList, DeviceResponse{
				ID:           dvc.ID,
				Number:       device.ID.User,
				PushName:     device.PushName,
				BusinessName: device.BusinessName,
				Contacts:     len(allContacts),
				Timestamp:    dvc.CreatedAt,
			})
		}
	}

	return deviceList, nil
}

// LogoutDeviceByJID logs out and removes a device from the store by its JID.
func LogoutDeviceByJID(jid string) {
	device, err := store.GetDeviceByJID(jid)
	if err != nil {
		return
	}
	err = store.RemoveDevice(device.ID)
	if err != nil {
		handler.FailOnError(err, "Error removing device from the store")
	}
}

// CheckIfNumberExistsAndGetJID verifies if a phone number is registered in WhatsApp and returns its JID.
func CheckIfNumberExistsAndGetJID(number string, client *whatsmeow.Client) (types.JID, error) {
	resp, err := client.IsOnWhatsApp([]string{number})
	if err != nil {
		return types.JID{}, fmt.Errorf("error checking if number exists: %v", err)
	}
	if resp[0].IsIn {
		return resp[0].JID, nil
	}
	return types.JID{}, fmt.Errorf("number is not registered in WhatsApp: %v", number)
}

// GenerateQRCode generates a QR code in base64-encoded format.
func GenerateQRCode(qrCode string) (string, error) {
	png, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to encode QR code: %w", err)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.Write(png); err != nil {
		return "", fmt.Errorf("failed to write PNG to buffer: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// GetClientInfo retrieves detailed information about a device and its contacts.
func GetClientInfo(deviceID int, client *whatsmeow.Client) DeviceInfoResponse {
	clientJID := types.NewJID(client.Store.ID.User, types.DefaultUserServer)

	// Get the profile picture of the client.
	picInfo, _ := client.GetProfilePictureInfo(clientJID, nil)
	var picURL string
	if picInfo != nil {
		picURL = picInfo.URL
	}

	// Get all contacts of the client.
	contacts, _ := client.Store.Contacts.GetAllContacts()
	var contactsResponse []ContactsResponse

	var wg sync.WaitGroup
	var mu sync.Mutex // Protects concurrent access to contactsResponse

	// Function to process each contact and add it to the response.
	processContact := func(key types.JID, contact types.ContactInfo) {
		defer wg.Done()

		contactName := contact.PushName
		if contactName == "" {
			contactName = contact.BusinessName
		}

		picContact, _ := client.GetProfilePictureInfo(key, nil)
		var contactProfileURL string
		if picContact != nil {
			contactProfileURL = picContact.URL
		}

		mu.Lock()
		contactsResponse = append(contactsResponse, ContactsResponse{
			Name:           contactName,
			Number:         key.User,
			ProfilePicture: contactProfileURL,
		})
		mu.Unlock()
	}

	for key, contact := range contacts {
		wg.Add(1)
		go processContact(key, contact)
	}

	wg.Wait()

	// Retrieve the active webhook URL for the device.
	var webhookURL string
	if webhook, err := GetWebhookActiveByDeviceID(deviceID); err == nil {
		webhookURL = webhook.WebhookURL
	}

	// Return device information response.
	return DeviceInfoResponse{
		DeviceID:         deviceID,
		ProfilePicture:   picURL,
		Webhook:          webhookURL,
		PhoneNumber:      clientJID.User,
		PushName:         client.Store.PushName,
		BusinessName:     client.Store.BusinessName,
		ContactsResponse: contactsResponse,
	}
}
