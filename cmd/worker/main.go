package main

import (
	"shin/src/app/workers"
	"shin/src/config"
	"time"

	"github.com/socious-io/gomail"
	"github.com/socious-io/gomq"
	database "github.com/socious-io/pkg_database"
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

	gomq.Setup(gomq.Config{
		Url:        config.Config.Nats.Url,
		Token:      config.Config.Nats.Token,
		ChannelDir: "shin",
		Consumers:  map[string]func(interface{}){},
	})

	//Initializing GoMail Library and Add it as Worker
	gomail.Setup(gomail.Config{
		ApiKey:         config.Config.Sendgrid.ApiKey,
		Url:            config.Config.Sendgrid.URL,
		DefaultFrom:    "info@socious.io",
		DefaultSubject: "Socious Verify",
		Templates:      config.Config.Sendgrid.Templates,
		WorkerChannel:  "shin/email",
		MessageQueue:   gomq.Mq,
	})

	workers.RegisterConsumers()

	gomq.Init()
}
