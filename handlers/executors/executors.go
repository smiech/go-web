package executorHandlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/smiech/go-web/globals"

	"github.com/fsnotify/fsnotify"
	models "github.com/smiech/go-web/models"
)

var List = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	v := globals.NetworkInfo
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

func ExecuteCommand(path string, output chan<- string, quit <-chan bool) {
	cmd := exec.Command("./scripts/echo.sh")

	cmd.Stdout = nil

	err := cmd.Start()
	go func() {
		select {
		case <-quit:
			log.Println("Killing process")
			cmd.Process.Kill()
		}
	}()

	if err != nil {
		log.Printf("Error!! : %v", err)
	}
}

func NewWatcher(path string, output chan<- string, quit <-chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				log.Println("Quiting filewatcher")
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					b, err := ioutil.ReadFile(event.Name) // just pass the file name
					if err != nil {
						fmt.Print(err)
					}

					//fmt.Println(b) // print the content as 'bytes'

					str := string(b) // convert content to a 'string'

					output <- str // print the content as a 'string'

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
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
