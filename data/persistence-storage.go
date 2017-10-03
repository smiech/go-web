package persistencestorage

import (
	"log"
	"time"

	models "github.com/adam72m/go-web/models"
	scribble "github.com/nanobox-io/golang-scribble"
)

type persistencestorage interface {
	getDevices(userID int) []models.Device
	addDevice(device models.Device)
	getUsers() []models.User
	init()
	storeHeartBeat(deviceID string, timestamp time.Time)
	getDeviceAlive(deviceID string) bool
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

type storageImplementation struct {
	DB *scribble.Driver
}

func (s storageImplementation) Init() {
	db, err := scribble.New("./db", nil)
	s.DB = db
	if err != nil {
		log.Println("Error", err)
	}
}

func (s storageImplementation) storeHeartBeat(deviceID string, timestamp time.Time) {
	devices, _ := s.DB.ReadAll("device")
	theDevice := first(devices, func(m models.Device) bool {
		return m.Id == deviceID
	})
	log.Printf
}

func (s storageImplementation) getDeviceAlive(deviceID string) bool {
	return true
}

func (s storageImplementation) getDevices(userID int) []models.Device {
	return nil
}

func (s storageImplementation) addDevice() {

}

func (s storageImplementation) getUsers() []models.User {
	return nil
}

// func (p *Device) save() error {
// 	res1B, _ := json.Marshal(p)
// 	fmt.Println(string(res1B))
// 	filename := "devices.json"
// 	return ioutil.WriteFile(filename, res1B, 0600)
// }

// func load(title string) (*Device, error) {
// 	filename := "devices.json"
// 	device := Device{}
// 	body, err := ioutil.ReadFile(filename)
// 	json.Unmarshal(body, &device)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Device{Id: device.Id, Name: device.Name}, nil
// }
