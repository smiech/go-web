package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

type Device struct {
	Id   string
	Name string
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

type Claims struct {
	Username string `json:"username"`
	// recommended having
	jwt.StandardClaims
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
