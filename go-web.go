package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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
		session, _ := store.Get(r, "session-name")
		r.ParseForm()
		usrName := r.Form["user"][0]
		currentTime := time.Now().String()
		session.Values[usrName] = currentTime
		session.Values["user"] = usrName
		session.Save(r, w)
		http.Redirect(w, r, "/step", http.StatusFound)
	}
}

func step(w http.ResponseWriter, r *http.Request) {
	wizardStep := &WizardStep{StepTitle: "Signup", StepNumber: "2", Content: "Register user"}
	renderTemplate(w, "wizardstep.tmpl", wizardStep)
}

func step2(w http.ResponseWriter, r *http.Request) {
	wizardStep := &WizardStep{StepTitle: "Setup", StepNumber: "3", Content: "Add devices"}
	renderTemplate(w, "setup.tmpl", wizardStep)
}

func step3(w http.ResponseWriter, r *http.Request) {
	wizardStep := &WizardStep{StepTitle: "Test", StepNumber: "4", Content: "Add devices"}
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

func Auth((w http.ResponseWriter, r *http.Request){
	authenticateResponse := AuthenticateResponse{}
}

func Authenticate(username string, password string) AuthenticateResponse {

	authenticateResponse := AuthenticateResponse{}
}

/*public AuthenticateResponse Authenticate(string username, string password)
  {
      AuthenticateResponse authenticateResponse = new AuthenticateResponse();

      var client = new RestClient(_settings.ServiceUrl);
      var request = new RestRequest("Account/AppLogin", Method.POST);
      request.AddHeader("Content-Type", "application/x-www-form-urlencoded");
      request.AddParameter("email", username);
      request.AddParameter("password", password);
      request.AddHeader("X-Requested-With", "XMLHttpRequest");
      request.AddHeader("nl.72media.riskmanager.appversion", "");

      RestResponse response = null;
      _log.LogInformation("Service starting task");
      Task.Run(async () =>
      {
          response = await RestClientHelper.GetResponseContentAsync(client, request) as RestResponse;
          var result = JsonConvert.DeserializeObject<AppLoginResponse>(response.Content);
          authenticateResponse.IsSuccessful = result.Success;

          if (result.Success)
          {
              authenticateResponse.User = new User
              {
                  ID = result?.UserId ?? Guid.Empty,
                  Email = result?.User?.Email ?? string.Empty,
                  ServiceCookies = new Dictionary<string, string>(),
                  Name = result?.User?.Name ?? string.Empty
              };

              response.Cookies.ToList().ForEach(x =>
              {
                  authenticateResponse.User.ServiceCookies.Add(x.Name, x.Value);
              });
          }

      }).Wait();

      return authenticateResponse;
  }*/

func main() {
	//configureLogger()
	loadConfiguration()
	loadTemplates()
	de := Device{Id: "aa", Name: "bb"}
	de.save()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", http.FileServer(http.Dir("./client/dist")))
	http.HandleFunc("/api/step", sessionCheckingHandler(step))
	http.HandleFunc("/api/step2", sessionCheckingHandler(step2))
	http.HandleFunc("/api/step3", sessionCheckingHandler(step3))
	http.HandleFunc("/api/login", login)
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
