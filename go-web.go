package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/oxtoacart/bpool"
)

var templates map[string]*template.Template
var bufpool *bpool.BufferPool

type WizardStep struct {
	StepTitle  string
	StepNumber string
	Content    string
}

type TemplateConfig struct {
	TemplateLayoutPath  string
	TemplateIncludePath string
}

var store = sessions.NewCookieStore([]byte("something-very-secret"))

var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`

var templateConfig TemplateConfig

func loadConfiguration() {
	templateConfig.TemplateLayoutPath = "templates/layouts/"
	templateConfig.TemplateIncludePath = "templates/"
}

func loadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layoutFiles, err := filepath.Glob(templateConfig.TemplateLayoutPath + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	includeFiles, err := filepath.Glob(templateConfig.TemplateIncludePath + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	mainTemplate := template.New("main")

	mainTemplate, err = mainTemplate.Parse(mainTmpl)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)
		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}

	log.Println("templates loading successful")

	bufpool = bpool.NewBufferPool(64)
	log.Println("buffer allocation successful")
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, fmt.Sprintf("The template %s does not exist.", name),
			http.StatusInternalServerError)
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

func index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.tmpl", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renderTemplate(w, "login.tmpl", nil)
	case "POST":
		r.ParseForm()
		log.Println(r.Form)
		session, _ := store.Get(r, "session-name")
		usrName := r.Form["user"][0]
		currentTime := time.Now().String()
		session.Values[usrName] = currentTime
		session.Values["user"] = usrName
		session.Save(r, w)
		http.Redirect(w, r, "/step", http.StatusFound)
	}
}

func step(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	// Set some session values.
	log.Println(session.Values["user"])
	log.Println(r.Form["user"])
	wizardStep := &WizardStep{StepTitle: "Signup", StepNumber: "1", Content: "Indian"}
	renderTemplate(w, "wizardstep.tmpl", wizardStep)
}

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
		log.Println(session.Values["user"])
		log.Println(r.Form["user"])
		fn(w, r)
	}
}

func main() {
	//configureLogger()
	loadConfiguration()
	loadTemplates()
	/*server := http.Server{
		Addr:    "127.0.0.1:" + "8080",
		Handler: context.ClearHandler(http.DefaultServeMux),
	}*/

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", login)
	http.HandleFunc("/step", step)
	http.HandleFunc("/login", sessionCheckingHandler(login))
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
