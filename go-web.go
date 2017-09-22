package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

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

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func configureLogger() {
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	defer f.Close()
	log.SetOutput(f)
	log.Println("hello")
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

type AuthenticateResponse struct {
	IsSuccessful bool
	User         User
}

func tokenAuth(w http.ResponseWriter, r *http.Request) {
	log.Printf("accessed \n")
}

func main() {
	configureLogger()

	r := mux.NewRouter()

	de := Device{Id: "aa", Name: "bb"}
	de.save()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./client/dist/")))
	r.Handle("/api/tokenauth", NotImplemented).Methods("POST")
	http.ListenAndServe(":8080", r)
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})
