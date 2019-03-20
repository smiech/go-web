package executorHandlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/smiech/go-web/handlers/recon"

	models "github.com/smiech/go-web/models"
)

var List = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	v := recon.NetworkInfo
	payload, _ := json.Marshal(v)
	log.Printf("List command called")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
	//msgTime := time.Unix(jsonMap.Time, 0)

})

var ExecuteHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	var jsonMap models.ExecuteData
	json.Unmarshal(body, &jsonMap)
	//msgTime := time.Unix(jsonMap.Time, 0)
	log.Printf("Execute command: %v Data: %v", jsonMap.CommandId, jsonMap.Data)
	if err != nil {
		log.Printf("Error parsing data %v ", err)
	}

	switch jsonMap.CommandId {
	case "start":
		handleStart(jsonMap.Data, w)
	case "stop":
		handleStop(jsonMap.Data, w)
	default:
		log.Printf("Not supported command: %v", jsonMap.CommandId)
	}
})

func handleStart(data string, w http.ResponseWriter) {
	log.Printf("Handling start command with data: %v", data)

	payload := `{
		"` + fmt.Sprintf("%v", data) + `"
	}`
	log.Printf("return payload: %v", payload)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(payload))
}

func handleStop(data string, w http.ResponseWriter) {
	log.Printf("Handling stop command with data: %v", data)
	payload := `{
			"` + fmt.Sprintf("%v", data) + `"
		}`
	log.Printf("return payload: %v", payload)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(payload))

}
