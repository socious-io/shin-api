package main

import (
	"shin/src/app"
	"shin/src/config"
	"shin/src/lib"
	"shin/src/services"
	"time"

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

	services.Connect()

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

	app.Serve()
}
