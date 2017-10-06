package deviceHandlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	persistence "github.com/adam72m/go-web/data"
	models "github.com/adam72m/go-web/models"
	"github.com/gorilla/mux"
)

var Storage persistence.PersistenceStorage

const aliveMarker string = "alive"
const alertMarker string = "alarm"
const stopRequestMarker string = "stop"

const stopCommand string = "FF000000000000"
const startCommand string = "000000000000FF"

var DeviceCallHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	var jsonMap models.CallbackData
	json.Unmarshal(body, &jsonMap)
	textDataBytes, _ := hex.DecodeString(jsonMap.Data)
	textData := fmt.Sprintf("%s", textDataBytes)
	msgTime := time.Unix(jsonMap.Time, 0)
	log.Printf("Device: %v Time: %v Data: %v", jsonMap.Device, msgTime, textData)
	if err != nil {
		log.Printf("Error parsing data %v ", err)
	}

	switch textData {
	case aliveMarker:
		handleHeartBeat(jsonMap.Device, jsonMap.Time)
	case alertMarker:
		log.Printf("Alert from device: %v", jsonMap.Device)
		handleAlertRequest(jsonMap.Device)
	case stopRequestMarker:
		log.Printf("Stop request from device: %v", jsonMap.Device)
		handleStopRequest(jsonMap.Device, w)
	default:
		log.Printf("Not supported command: %v from device: %v", textData, jsonMap.Device)
	}
})

func handleAlertRequest(deviceId string) {
	device, _ := Storage.GetDeviceByGuid(deviceId)
	Storage.RegisterDeviceEvent(device.Id, models.DeviceEvent{Name: "alert", CreateTime: time.Now(), DeviceId: device.Id})
}

func handleStopRequest(deviceId string, w http.ResponseWriter) {
	device, _ := Storage.GetDeviceByGuid(deviceId)
	Storage.RegisterDeviceEvent(device.Id, models.DeviceEvent{Name: stopRequestMarker, CreateTime: time.Now(), DeviceId: device.Id})

	deviceCommands, _ := Storage.GetDeviceCommands(device.Id)
	var command models.DeviceCommand
	timeRef := time.Time{}
	for _, v := range deviceCommands {
		if v.ExecutedTime == timeRef {
			command = v
			break

		}
	}
	if command != (models.DeviceCommand{}) {
		var commandToSend string
		if command.Command == "stop" {
			log.Printf("Sending downlink command to stop for device: %v Command: %v", deviceId, commandToSend)
			commandToSend = stopCommand
		} else {
			log.Printf("Sending downlink command to start for device: %v Command: %v", deviceId, commandToSend)
			commandToSend = startCommand
		}
		payload := `{
			"` + fmt.Sprintf("%v", deviceId) + `" : { "downlinkData" : "` + commandToSend + `" }
		}`
		log.Printf("return payload: %v", payload)
		w.Header().Set("content-type", "application/json")
		w.Write([]byte(payload))
	} else {
		log.Printf("No command to execute for device %v", deviceId)
	}
}

func handleHeartBeat(deviceId string, timeEpoch int64) {
	devices, err := Storage.GetDevices(0)
	log.Printf("%v", err)
	timeStamp := time.Unix(timeEpoch, 0)
	if len(devices) == 0 {
		Storage.AddDevice(models.Device{Id: 1, Guid: deviceId, Name: "Adam test 1", LastSeen: time.Unix(timeEpoch, 0)})
	}

	log.Printf("Existing device timestamp: %v", devices[0].LastSeen)
	Storage.StoreHeartBeat(deviceId, timeStamp)
}

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entered the status handler")
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	isDeviceAlive, _ := Storage.GetDeviceAlive(deviceId)
	var response string
	if isDeviceAlive {
		response = fmt.Sprint("ALIVE")
	} else {
		response = fmt.Sprint("NOTALIVE")
	}
	w.Write([]byte(response))
})

var GetDevicesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	v, _ := Storage.GetDevices(0)
	payload, _ := json.Marshal(v)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})
