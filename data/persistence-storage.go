package persistencestorage

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	models "github.com/adam72m/go-web/models"
	scribble "github.com/nanobox-io/golang-scribble"
)

const deviceTable string = "device"
const userTable string = "user"
const deviceEventTable string = "deviceEvent"
const deviceCommandTable string = "deviceCommand"

type PersistenceStorage interface {
	GetDevices(userID int) ([]models.Device, error)
	AddDevice(device models.Device) error
	GetUsers() ([]models.User, error)
	AddUser(user models.User) error
	StoreHeartBeat(deviceID string, timestamp time.Time) error
	GetDeviceAlive(deviceID string) (bool, error)
	RegisterDeviceEvent(deviceId int, event models.DeviceEvent) error
	RegisterDeviceCommand(deviceId int, command models.DeviceCommand) error
	GetDeviceCommands(deviceId int) ([]models.DeviceCommand, error)
}

func first(vs []models.Device, f func(models.Device) bool) models.Device {
	for _, v := range vs {
		if f(v) {
			return v
		}
	}
	log.Fatal("No elements found")
	return models.Device{}
}

type StorageError struct {
	Message string
}

func (s StorageError) Error() string {
	return fmt.Sprintf("Storage error: %v", s.Message)
}

type StorageImplementation struct {
	DB *scribble.Driver
}

func (s StorageImplementation) StoreHeartBeat(deviceID string, timestamp time.Time) error {
	devs, err := s.GetDevices(0)
	if err != nil {
		return StorageError{Message: "No Devices found"}
	}
	theDevice := first(devs, func(m models.Device) bool {
		return m.Guid == deviceID
	})
	theDevice.LastSeen = time.Now()
	s.DB.Write("device", fmt.Sprintf("%v", theDevice.Id), theDevice)
	log.Printf("%v", theDevice)
	return nil
}

func (s StorageImplementation) GetDeviceAlive(deviceID string) (bool, error) {
	device := models.Device{}
	err := s.DB.Read(deviceTable, deviceID, &device)
	if err != nil {
		return false, StorageError{Message: "No device found"}
	}
	isAlive := time.Since(device.LastSeen).Hours() < 1
	return isAlive, nil
}

func (s StorageImplementation) GetDevices(userID int) ([]models.Device, error) {
	devices, err := s.DB.ReadAll(deviceTable)
	if err != nil {
		return nil, StorageError{Message: "Failed to read all devices"}
	}
	devs := []models.Device{}
	for _, devic := range devices {
		f := models.Device{}
		json.Unmarshal([]byte(devic), &f)
		devs = append(devs, f)
	}
	return devs, nil
}

func (s StorageImplementation) AddDevice(device models.Device) error {
	err := s.DB.Write(deviceTable, fmt.Sprintf("%v", device.Id), device)
	if err != nil {
		return StorageError{Message: "Failed adding device"}
	}
	log.Printf("Device added: %v")
	return nil
}

func (s StorageImplementation) GetUsers() ([]models.User, error) {
	users, err := s.DB.ReadAll(userTable)
	if err != nil {
		return nil, StorageError{Message: "Failed to read all users"}
	}
	usrs := []models.User{}
	for _, u := range users {
		f := models.User{}
		json.Unmarshal([]byte(u), &f)
		usrs = append(usrs, f)
	}
	return usrs, nil
}

func (s StorageImplementation) AddUser(user models.User) error {
	err := s.DB.Write(userTable, fmt.Sprintf("%v", user.Id), user)
	if err != nil {
		return StorageError{Message: "Failed adding user"}
	}
	log.Printf("User added: %v", user)
	return nil
}

func (s StorageImplementation) RegisterDeviceEvent(deviceId int, event models.DeviceEvent) error {
	currentDeviceEvents, err := s.DB.ReadAll(deviceEventTable)
	if err != nil {
		return StorageError{Message: "Failed registering event"}
	}
	newId := len(currentDeviceEvents)
	event.Id = newId
	err = s.DB.Write(deviceEventTable, fmt.Sprintf("%v", event.Id), event)
	if err != nil {
		return StorageError{Message: "Failed registering event"}
	}
	log.Printf("event added: %v", event)
	return nil
}

func (s StorageImplementation) RegisterDeviceCommand(deviceId int, command models.DeviceCommand) error {
	currentDeviceCommands, err := s.DB.ReadAll(deviceCommandTable)
	if err != nil {
		return StorageError{Message: "Failed registering command"}
	}
	newId := len(currentDeviceCommands)
	command.Id = newId
	err = s.DB.Write(deviceCommandTable, fmt.Sprintf("%v", command.Id), command)
	if err != nil {
		return StorageError{Message: "Failed registering command"}
	}
	log.Printf("command added: %v", command)
	return nil
}

func (s StorageImplementation) GetDeviceCommands(deviceId int) ([]models.DeviceCommand, error) {
	deviceCommandsArray, err := s.DB.ReadAll(deviceCommandTable)
	if err != nil {
		return nil, StorageError{Message: "Failed to read all commands"}
	}
	dvs := []models.DeviceCommand{}
	for _, u := range deviceCommandsArray {
		f := models.DeviceCommand{}
		json.Unmarshal([]byte(u), &f)
		dvs = append(dvs, f)
	}
	return dvs, nil
}
