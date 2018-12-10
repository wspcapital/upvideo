package main

import (
	"bitbucket.org/marketingx/upvideo/app/infrastructure/web"
)

type SessionConfig struct {
	Storage    string
	TTLMinutes int32
}

type Config struct {
	Dsn               string
	CheckInterval     int
	Session SessionConfig
	WebServer web.WebServerParams
}
