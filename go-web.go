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

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const username string = "adam"
const password string = "enter"
const wwwRoot = "./client/dist/"

// var domain = "localhost"
// var httpPort = "8080"
// var httpsPort = "8443"

var port string = "8080"

type Device struct {
	Id   string
	Name string
}

func (p *Device) save() error {
	res1B, _ := json.Marshal(p)
	fmt.Println(string(res1B))
	filename := "devices.json"
	return ioutil.WriteFile(filename, res1B, 0600)
}

func load(title string) (*Device, error) {
	filename := "devices.json"
	device := Device{}
	body, err := ioutil.ReadFile(filename)
	json.Unmarshal(body, &device)
	if err != nil {
		return nil, err
	}
	return &Device{Id: device.Id, Name: device.Name}, nil
}

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

func sessionCheckingHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, fmt.Sprintf("The session expired \n"), http.StatusInternalServerError)
			return
		}
		userName := session.Values["user"]
		log.Printf("%v", userName)
		if userName != "adam" {
			http.Error(w, fmt.Sprintf("Invalid login \n"), http.StatusInternalServerError)
			return
		}

		fn(w, r)
	}
}

type User struct {
	Guid           string
	Name           string
	Email          string
	ServiceCookies map[string]string
}

type Credentials struct {
	Username string
	Password string
}

type AuthenticateResponse struct {
	IsSuccessful bool
	User         User
}

func tokenAuth(w http.ResponseWriter, r *http.Request) {
	log.Printf("accessed \n")
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
	//de := Device{Id: "aa", Name: "bb"}
	//de.save()
	r.Handle("/api/v1/login", loginHandler).Methods("POST", "OPTIONS")
	//r.Handle("/api/v1/login", optionHandler).Methods("OPTIONS")
	r.Handle("/api/v1/status", authMiddleware(ProductsHandler)).Methods("GET")

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

type Claims struct {
	Username string `json:"username"`
	// recommended having
	jwt.StandardClaims
}

func GetCredentials(reqBody io.ReadCloser) Credentials {
	decoder := json.NewDecoder(reqBody)
	var t Credentials
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
	//expireCookie := time.Now().Add(time.Hour * 1)

	/* Create a map to store our claims*/
	claims := Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:" + port,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)
	// Place the token in the client's cookie
	//cookie := http.Cookie{Name: "Auth", Value: tokenString, Expires: expireCookie, HttpOnly: false}
	resp := AuthRequestResult{
		State: 1,
		Data:  DataStruct{AccessToken: tokenString, User: User{Name: "Adam"}},
	}
	payload, err := json.Marshal(resp)
	log.Printf("%v", err)
	check := AuthRequestResult{}
	json.Unmarshal(payload, &check)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))

	// AuthRequestResult {
	// 	State: number;
	// 	Msg: string;
	// 	Data: Object;
	// }

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
		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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

type DataStruct struct {
	AccessToken string `json:"accessToken"`
	User        User   `json:"user"`
}

type AuthRequestResult struct {
	State int
	Msg   string
	Data  DataStruct
}

type Product struct {
	Id          int
	Name        string
	Slug        string
	Description string
}

/* We will create our catalog of VR experiences and store them in a slice. */
var products = []Product{
	Product{Id: 1, Name: "Hover Shooters", Slug: "hover-shooters", Description: "Shoot your way to the top on 14 different hoverboards"},
	Product{Id: 2, Name: "Ocean Explorer", Slug: "ocean-explorer", Description: "Explore the depths of the sea in this one of a kind underwater experience"},
	Product{Id: 3, Name: "Dinosaur Park", Slug: "dinosaur-park", Description: "Go back 65 million years in the past and ride a T-Rex"},
	Product{Id: 4, Name: "Cars VR", Slug: "cars-vr", Description: "Get behind the wheel of the fastest cars in the world."},
	Product{Id: 5, Name: "Robin Hood", Slug: "robin-hood", Description: "Pick up the bow and arrow and master the art of archery"},
	Product{Id: 6, Name: "Real World VR", Slug: "real-world-vr", Description: "Explore the seven wonders of the world in VR"},
}

/* The status handler will be invoked when the user calls the /status route
It will simply return a string with the message "API is up and running" */
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running"))
})

/* The products handler will be called when the user makes a GET request to the /products endpoint.
This handler will return a list of products available for users to review */
var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Here we are converting the slice of products to json
	payload, _ := json.Marshal(products)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})
