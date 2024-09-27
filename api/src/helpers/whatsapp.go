package helpers

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
	"whatsgoingon/store"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ztrue/tracerr"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Códigos de erro
var (
	dbLog                 = waLog.Stdout("Database", "WARN", true)
	wmLog                 = waLog.Stdout("WhatsMeow", "WARN", true)
	container             *sqlstore.Container
	client                *whatsmeow.Client
	once                  sync.Once
	dbMutex               sync.Mutex
	ErrDBConnectionFailed = errors.New("failed to connect to the database")
	ErrDeviceNotFound     = errors.New("device not found in the store")
	ErrClientConnection   = errors.New("failed to connect the WhatsApp client")
	ErrMessageSending     = errors.New("failed to send the message")
)

type DeviceResponse struct {
	ID           string `json:"id"`
	Number       string `json:"number"`
	PushName     string `json:"push_name"`
	BusinessName string `json:"business_name"`
	Contacts     int    `json:"contacts"`
	Timestamp	 time.Time `json:"timestamp"`
}

// ConnectToDatabase connects to the database.
func connectToDatabase() (*sqlstore.Container, error) {
	dbUser := os.Getenv("pg_username")
	dbPwd := os.Getenv("pg_password")
	dbTCPHost := os.Getenv("pg_hostname")
	dbPort := os.Getenv("pg_port")
	dbName := os.Getenv("pg_database")

	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s sslmode=disable",
		dbTCPHost, dbUser, dbPwd, dbPort, dbName)

	container, err := sqlstore.New("pgx", dbURI, dbLog)
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrDBConnectionFailed, err))
	}
	return container, nil
}

// GetClient returns a WhatsApp client.
func GetClient() (*whatsmeow.Client, error) {
	waLog.Stdout("WhatsappHelper", "WARN", true)
	var err error
	if client == nil {
		once.Do(func() {
			dbMutex.Lock()
			defer dbMutex.Unlock()

			container, err = connectToDatabase()
			if err != nil {
				return // Error already wrapped
			}

			deviceStore, err := container.GetFirstDevice()
			if err != nil {
				err = tracerr.Wrap(fmt.Errorf("%w: %v", ErrDeviceNotFound, err))
				return
			}
			client = whatsmeow.NewClient(deviceStore, wmLog)
		})
	}
	if err != nil {
		return nil, err
	}
	if !client.IsConnected() {
		err := client.Connect()
		if err != nil {
			return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrClientConnection, err))
		}
	}
	return client, nil
}

// GetWhatsAppClientByJID returns a WhatsApp client by JID.
func GetWhatsAppClientByJID(whatsappID string) (*whatsmeow.Client, error) {
	waLog.Stdout("WhatsappHelper", "WARN", true)
	
	dbMutex.Lock()
	defer dbMutex.Unlock()

	container, _ = connectToDatabase()

	jid, _ := types.ParseJID(whatsappID)
	
	deviceStore, err := container.GetDevice(jid)
	if err != nil {
		err = tracerr.Wrap(fmt.Errorf("1-----------> %w: %v", ErrDeviceNotFound, err))
		return nil, err
	}

	client := whatsmeow.NewClient(deviceStore, wmLog)
	if client.Store.ID == nil {
		return nil, tracerr.Wrap(fmt.Errorf("2--------> %w: %v", ErrClientConnection, client.Store.ID))
	}
	
	if !client.IsConnected() {
		err := client.Connect()
		if err != nil {
			return nil, tracerr.Wrap(fmt.Errorf("3------> %w: %v", ErrClientConnection, err))
		}
	}
	return client, nil
}

// GetWhatsappClientByDeviceID returns a WhatsApp client by device ID.
func GetWhatsappClientByDeviceID(deviceID string) (*whatsmeow.Client, error) {
	device, _ := store.GetDeviceById(deviceID)
	client, err := GetWhatsAppClientByJID(device.JID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetAllWhatsappIDs() ([]string, error) {
	container, err := connectToDatabase()
	if err != nil {
		return nil, err // Error already wrapped
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

func NewClient() (*whatsmeow.Client, error) {
	container, err := connectToDatabase()
	if err != nil {
		return nil, err // Error already wrapped
	}

	deviceStore := container.NewDevice()
	client := whatsmeow.NewClient(deviceStore, wmLog)

	return client, nil
}

func GetDeviceList() ([]DeviceResponse, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	container, err := connectToDatabase()
	if err != nil {
		return nil, err // Error already wrapped
	}

	devices, err := container.GetAllDevices()
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrDeviceNotFound, err))
	}

	var deviceList []DeviceResponse
	for _, device := range devices {
		allContacts, err := device.Contacts.GetAllContacts()
		if err != nil {
			// Handle error gracefully, perhaps log it or return an error
			continue // Skip this device and continue with the next one
		}

		if dvc, err := store.GetDeviceByJID(device.ID.String()); err != nil {
			continue
		} else {
			deviceList = append(deviceList, DeviceResponse{
				ID:           dvc.DeviceID(),
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