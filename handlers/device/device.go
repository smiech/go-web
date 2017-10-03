package deviceHandlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	models "github.com/adam72m/go-web/models"
)

var DataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	/* Create the token */
	log.Printf("data handler2 (downlink) invoked")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("r.PostForm", r.PostForm)
	log.Println("r.Form", r.Form)

	body, err := ioutil.ReadAll(r.Body)
	var jsonMap map[string]interface{}
	json.Unmarshal(body, &jsonMap)
	log.Printf("Body: %v", body)
	if err != nil {
		log.Printf("Error parsing data %v ", err)
	}
	payload := `{
		"` + fmt.Sprintf("%v", jsonMap["device"]) + `" : { "downlinkData" : "deadbeefcafebabe"}
	}`

	log.Printf("return payload: %v", payload)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(payload))
})

/* We will create our catalog of VR experiences and store them in a slice. */
var devices = []models.Device{
	models.Device{Id: "004d2BBB", Name: "Adam test 1"},
}

/* The status handler will be invoked when the user calls the /status route
It will simply return a string with the message "API is up and running" */
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running"))
})

/* The products handler will be called when the user makes a GET request to the /products endpoint.
This handler will return a list of products available for users to review */
var GetDevicesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Here we are converting the slice of products to json
	payload, _ := json.Marshal(devices)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})
