package main

import (
	"log"
	"shin/src/app"
	"shin/src/config"
	"shin/src/lib"
	"time"

	"github.com/socious-io/goaccount"
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

	if config.Config.Storage.Type == "AWS" {
		lib.InitS3Lib(lib.S3ConfigType{
			AccessKeyId:     config.Config.Storage.S3.AccessKeyId,
			SecretAccessKey: config.Config.Storage.S3.SecretAccessKey,
			DefaultRegion:   config.Config.Storage.S3.DefaultRegion,
			Bucket:          config.Config.Storage.S3.Bucket,
			CDNUrl:          config.Config.Storage.S3.CDNUrl,
		})
	} else if config.Config.Storage.Type == "GCS" {
		lib.InitGCSLib(lib.GCSConfigType{
			Bucket:          config.Config.Storage.GCS.Bucket,
			CDNUrl:          config.Config.Storage.GCS.CDNUrl,
			CredentialsPath: config.Config.Storage.GCS.CredentialsPath,
		})
	}

	//Initializing GoMQ Library
	gomq.Setup(gomq.Config{
		Url:        config.Config.Nats.Url,
		Token:      config.Config.Nats.Token,
		ChannelDir: "shin",
	})
	gomq.Connect()

	if err := goaccount.Setup(config.Config.GoAccounts); err != nil {
		log.Fatalf("goaccount error %v", err)
	}

	app.Serve()
}
