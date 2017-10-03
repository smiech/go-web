package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	deviceHandlers "github.com/adam72m/go-web/handlers/device"
	m "github.com/adam72m/go-web/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const username string = "adam"
const password string = "enter"
const wwwRoot = "./client/dist/"

var port string = "8080"

var store = sessions.NewCookieStore([]byte("dwadziescia-muharadzinow-bije-trzech-rabinow"))

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
	port = os.Getenv("HTTP_PLATFORM_PORT")
	if port == "" {
		port = os.Getenv("port")
		if port == "" {
			port = "8080"
		}

	}

	r := mux.NewRouter()
	r.Handle("/api/v1/login", loginHandler).Methods("POST", "OPTIONS")
	r.Handle("/api/v1/submit", deviceHandlers.DataHandler).Methods("POST", "OPTIONS")
	r.Handle("/api/v1/devices", authMiddleware(deviceHandlers.GetDevicesHandler)).Methods("GET")

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

var dataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	/* Create the token */
	log.Printf("data handler invoked")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("r.PostForm", r.PostForm)
	log.Println("r.Form", r.Form)

	body, err := ioutil.ReadAll(r.Body)
	log.Printf("Body: %v", body)
	if err != nil {
		log.Printf("Error parsing data %v ", err)
	}
	var jsonMap map[string]interface{}
	json.Unmarshal(body, &jsonMap)
	log.Printf("%v", jsonMap)
})

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If no Auth cookie is set then return a 404 not found
		cookie, err := r.Cookie("Auth")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Return a Token using the cookie
		token, err := jwt.ParseWithClaims(cookie.Value, &m.Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Make sure token's signature wasn't changed
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected siging method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
