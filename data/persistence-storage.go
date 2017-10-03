package persistencestorage

import (
	models "github.com/adam72m/go-web/models"
)

type persistencestorage interface {
	getDevices(userID int) []models.Device
	addDevice(device models.Device)
	getUsers() []models.User
}

type storageImplementation struct {
}

func (s storageImplementation) getDevices(userID int) []models.Device {
	return nil
}

func (s storageImplementation) addDevice() {
}

func (s storageImplementation) getUsers() []models.User {
	return nil
}
