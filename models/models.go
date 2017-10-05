package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type CallbackData struct {
	Device string
	Data   string
	Time   int64
}

type DeviceEvent struct {
	Id         int
	DeviceId   int
	Name       string
	CreateTime time.Time
}

type DeviceCommand struct {
	Id           int
	DeviceId     int
	Command      string
	CreateTime   time.Time
	ExecutedTime time.Time
}

type Device struct {
	Guid     string
	Id       int
	Name     string
	LastSeen time.Time
}

type User struct {
	Id             int
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
