package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	executeHandlers "github.com/smiech/go-web/handlers/executors"
	m "github.com/smiech/go-web/models"
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

	env := os.Args
	var fileWriter io.Writer
	envLength := len(env)
	if envLength > 0 && env[envLength-1] == "debug" || os.Getenv("debug") == "true" {
		fileWriter = os.Stdout
	} else {
		fileWriter = configureLogger()
	}
	executeHandlers.ExecuteCommand("./scripts/echo.sh")
	port = os.Getenv("HTTP_PLATFORM_PORT")
	if port == "" {
		port = os.Getenv("port")
		if port == "" {
			port = "8080"
		}

	}

	r := mux.NewRouter()
	r.Handle("/api/v1/login", loginHandler).Methods("POST", "OPTIONS")
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

func GetCredentials(reqBody io.ReadCloser) m.Credentials {
	decoder := json.NewDecoder(reqBody)
	var t m.Credentials
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	defer reqBody.Close()
	return t
}

/* Set up a global string for our secret */
var mySigningKey = []byte("secret")

/* Handlers */
var loginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	/* Create the token */
	log.Printf("login handler invoked")

	credentials := GetCredentials(r.Body)

	if credentials.Username != username && credentials.Password != password {
		w.WriteHeader(http.StatusForbidden)
		log.Printf("Error in request")
		return
	}
	expireToken := time.Now().Add(time.Hour * 1).Unix()

	/* Create a map to store our claims*/
	claims := m.Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:" + port,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(mySigningKey)

	resp := m.AuthRequestResult{
		State: 1,
		Data:  m.DataStruct{AccessToken: tokenString, User: m.User{Name: "Adam"}},
	}
	payload, err := json.Marshal(resp)
	log.Printf("%v", err)
	check := m.AuthRequestResult{}
	json.Unmarshal(payload, &check)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})
