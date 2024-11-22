package main

import (
	"shin/src/config"
	"shin/src/database"
	"shin/src/lib"
	"shin/src/services"
	"time"
)

func main() {

	config.Init("config.yml")
	database.Connect(&database.ConnectOption{
		URL:         config.Config.Database.URL,
		SqlDir:      config.Config.Database.SqlDir,
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
	})

	lib.InitSendGridLib(lib.SendGridType{
		Disabled: config.Config.Sendgrid.Disabled,
		ApiKey:   config.Config.Sendgrid.ApiKey,
		Url:      config.Config.Sendgrid.URL,
	})

	services.Init()
}
