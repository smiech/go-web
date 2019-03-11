package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	fileParser "github.com/smiech/go-web/helpers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	executeHandlers "github.com/smiech/go-web/handlers/executors"
)

const username string = "adam"
const password string = "enter"
const wwwRoot = "./client/dist/"

var port string = "8080"

func configureLogger() io.Writer {
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	log.SetOutput(f)
	log.Println("Logging configured")
	return f
}

func main() {
	output := make(chan string)
	quit := make(chan bool)
	quit2 := make(chan bool)
	env := os.Args
	var fileWriter io.Writer
	envLength := len(env)
	if envLength > 0 && env[envLength-1] == "debug" || os.Getenv("debug") == "true" {
		fileWriter = os.Stdout
	} else {
		fileWriter = configureLogger()
	}
	executeHandlers.ExecuteCommand("./scripts/echo.sh", output, quit)
	go executeHandlers.NewWatcher("./dumps", output, quit2)
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
			case <-time.After(50 * time.Second):
				log.Println("Sending quit signal")
				quit <- true
				quit2 <- true
			}
		}
	}()
	port = os.Getenv("HTTP_PLATFORM_PORT")
	if port == "" {
		port = os.Getenv("port")
		if port == "" {
			port = "8080"
		}

	}

	r := mux.NewRouter()
	r.Handle("/api/v1/execute", executeHandlers.ExecuteHandler).Methods("POST", "OPTIONS")

	r.HandleFunc("/", indexHandler)
	r.PathPrefix("/").HandlerFunc(staticHandler)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With",
		"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type",
	})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	corsHandler := handlers.CORS(originsOk, headersOk, methodsOk)(r)
	log.Printf("%v", http.ListenAndServe(":"+port, handlers.LoggingHandler(fileWriter, corsHandler)))

}

// Redirect all traffic to HTTPS
// func redirectHandler(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "https://"+domain+":"+httpsPort+r.RequestURI, http.StatusMovedPermanently)
// }

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, wwwRoot+"index.html")
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	fileSystemPath := wwwRoot + r.URL.Path
	endURIPath := strings.Split(requestPath, "/")[len(strings.Split(requestPath, "/"))-1]
	splitPath := strings.Split(endURIPath, ".")
	if len(splitPath) > 1 {
		if f, err := os.Stat(fileSystemPath); err == nil && !f.IsDir() {
			http.ServeFile(w, r, fileSystemPath)
			return
		}
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, wwwRoot+"index.html")
}

/* Set up a global string for our secret */
var mySigningKey = []byte("secret")
