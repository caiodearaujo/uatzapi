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

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/skip2/go-qrcode"
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
	dbMutex               sync.Mutex
	ErrDBConnectionFailed = errors.New("failed to connect to the database")
	ErrDeviceNotFound     = errors.New("device not found in the store")
	ErrClientConnection   = errors.New("failed to connect the WhatsApp client")
	ErrMessageSending     = errors.New("failed to send the message")
)

type DeviceResponse struct {
	ID           int       `json:"id"`
	Number       string    `json:"number"`
	PushName     string    `json:"push_name"`
	BusinessName string    `json:"business_name"`
	Contacts     int       `json:"contacts"`
	Timestamp    time.Time `json:"timestamp"`
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

// GetWhatsAppClientByJID returns a WhatsApp client by JID.
func GetWhatsAppClientByJID(whatsappID string) (*whatsmeow.Client, error) {
	waLog.Stdout("WhatsappHelper", "WARN", true)

	dbMutex.Lock()
	defer dbMutex.Unlock()

	container, _ = connectToDatabase()

	jid, _ := types.ParseJID(whatsappID)

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

// GetWhatsappClientByDeviceID returns a WhatsApp client by device ID.
func GetWhatsappClientByDeviceID(deviceID int) (*whatsmeow.Client, error) {
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

func LogoutDeviceByJID(jid string) {
	device, err := store.GetDeviceByJID(jid)
	if err != nil {
		return
	}
	err, _ = store.RemoveDevice(device.ID)
	if err != nil {
		handler.FailOnError(err, "Error removing device from the store")
	}
	return
}

func CheckIfNumberExistsAndGetJID(number string, client *whatsmeow.Client) (types.JID, error) {
	numberList := []string{number}
	resp, err := client.IsOnWhatsApp(numberList)
	if err != nil {
		return types.JID{}, fmt.Errorf("error checking if number exists: %v", err)
	}
	if resp[0].IsIn {
		return resp[0].JID, nil
	}
	return types.JID{}, fmt.Errorf("number is not registered in WhatsApp: %v", number)
}

func GenerateQRCode(qrCode string) (string, error) {
	// Gera a imagem do QR code
	png, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to encode QR code: %w", err)
	}

	// Converte o PNG para base64
	buf := new(bytes.Buffer)
	if _, err := buf.Write(png); err != nil { // Corrigido aqui
		return "", fmt.Errorf("failed to write PNG to buffer: %w", err)
	}
	qrCodeBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return qrCodeBase64, nil
}
