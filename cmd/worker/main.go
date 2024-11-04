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
	}, map[string]string{
		"otp":                    "d-5ace79d4674b45d1bfb9b35c4d6eb8c0",
		"forget-password":        "d-d68cc8d8409942599f761261e5a7fbcb",
		"credentials-recipients": "d-12fddb16345e4073a6741884237ed39c",
	})

	services.Init()
}
