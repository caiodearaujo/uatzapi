package helpers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ztrue/tracerr"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// Códigos de erro
var (
	dbLog                 = waLog.Stdout("Database", "DEBUG", true)
	wmLog                 = waLog.Stdout("WhatsMeow", "DEBUG", true)
	container             *sqlstore.Container
	client                *whatsmeow.Client
	once                  sync.Once
	dbMutex               sync.Mutex
	ErrDBConnectionFailed = errors.New("failed to connect to the database")
	ErrDeviceNotFound     = errors.New("device not found in the store")
	ErrClientConnection   = errors.New("failed to connect the WhatsApp client")
	ErrMessageSending     = errors.New("failed to send the message")
)

type Device struct {
	ID           string `json:"id"`
	Number       string `json:"number"`
	PushName     string `json:"push_name"`
	BusinessName string `json:"business_name"`
	Contacts     int    `json:"contacts"`
}

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

func GetClient() (*whatsmeow.Client, error) {
	waLog.Stdout("WhatsappHelper", "DEBUG", true)
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

func GetClientById(clientID string) (*whatsmeow.Client, error) {
	waLog.Stdout("WhatsappHelper", "DEBUG", true)
	var err error
	dbMutex.Lock()
	defer dbMutex.Unlock()

	container, _ = connectToDatabase()

	jid, _ := types.ParseJID(clientID)
	deviceStore, err := container.GetDevice(jid)

	if err != nil {
		err = tracerr.Wrap(fmt.Errorf("%w: %v", ErrDeviceNotFound, err))
		return nil, err
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

func GetAllClientIDs() ([]string, error) {
	container, err := connectToDatabase()
	if err != nil {
		return nil, err // Error already wrapped
	}

	deviceStore, err := container.GetAllDevices()
	if err != nil {
		return nil, tracerr.Wrap(fmt.Errorf("%w: %v", ErrDeviceNotFound, err))
	}

	var clientIDs []string

	for _, device := range deviceStore {
		clientIDs = append(clientIDs, device.ID.String())
	}

	return clientIDs, nil
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

func SendMessage(message string, recipient string) error {
	_, err := GetClient()
	if err != nil {
		return err
	}

	destination := types.JID{
		User:   recipient,
		Server: types.DefaultUserServer,
	}

	encryptedMessage := &waE2E.Message{
		Conversation: proto.String(message),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Timeout context
	defer cancel()

	_, err = client.SendMessage(ctx, destination, encryptedMessage)
	if err != nil {
		return tracerr.Wrap(fmt.Errorf("%w: %v", ErrMessageSending, err))
	}

	return nil
}

func GetDeviceList() ([]Device, error) {
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

	var deviceList []Device
	for _, device := range devices {
		allContacts, err := device.Contacts.GetAllContacts()
		if err != nil {
			// Handle error gracefully, perhaps log it or return an error
			continue // Skip this device and continue with the next one
		}

		deviceList = append(deviceList, Device{
			ID:           device.ID.String(),
			Number:       device.ID.User,
			PushName:     device.PushName,
			BusinessName: device.BusinessName,
			Contacts:     len(allContacts),
		})
	}

	return deviceList, nil
}
