package adminHandlers

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
)

var Storage persistence.PersistenceStorage

var AddCommandHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	default:
		log.Printf("Not supported command: %v from device: %v", textData, jsonMap.Device)
	}
})
