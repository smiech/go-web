package recon

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/go-cmd/cmd"
	fileParser "github.com/smiech/go-web/helpers"
	"github.com/smiech/go-web/models"
)

var isRunning = false
var isSpecificRunning = false
var isFileMonitorRunning = false
var scanCmd = cmd.NewCmdOptions(cmd.Options{Streaming: true}, "./scripts/all.sh")
var scanSpecificCmd *cmd.Cmd
var NetworkInfo models.NetworkInfo
var output = make(chan string)
var quit2 = make(chan bool)

func StartSpecific(bssid string) {
	if isRunning || isSpecificRunning {
		log.Println("Tried to start specific recon when it's already running!")
		return
	}
	log.Println("Starting recon on ", bssid)
	isSpecificRunning = true
	scanSpecificCmd = cmd.NewCmdOptions(cmd.Options{Streaming: true}, "./scripts/specific.sh", fmt.Sprintf("-b=%v", bssid))
	scanSpecificCmd.Stdout = nil
	scanSpecificCmd.Start() // non-blocking
	if !isFileMonitorRunning {
		startWatching()
	}
}

func Start() {
	if isRunning || isSpecificRunning {
		log.Println("Tried to start recon when it's already running!")
		return
	}
	log.Println("Starting recon")
	isRunning = true
	scanCmd.Stdout = nil
	scanCmd.Start() // non-blocking
	if !isFileMonitorRunning {
		startWatching()
	}
}

func Stop() {
	if !isRunning {
		log.Println("Tried to stop recon when it's already stopped!")
		return
	}
	log.Println("Stopping recon")
	scanCmd.Stop()
	quit2 <- true
	isRunning = false
	isFileMonitorRunning = false
}

func StopSpecific() {
	if !isSpecificRunning {
		log.Println("Tried to stop recon when it's already stopped!")
		return
	}
	log.Println("Stopping specific recon")
	scanSpecificCmd.Stop()
	quit2 <- true
	isSpecificRunning = false
	isFileMonitorRunning = false
}

func startWatching() {
	isFileMonitorRunning = true
	go newWatcher("./dumps", output, quit2)

	go func() {
		for {
			select {
			case file := <-output:
				log.Println("File contents:")
				log.Println("modified file:", file)
				records, err := fileParser.Parse(file)
				if err != nil {
					fmt.Printf("error opening file: %v", err)
				}
				log.Println(records)
				NetworkInfo = records
			}
		}
	}()
}

func newWatcher(path string, output chan<- string, quit <-chan bool) {
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
