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

	lib.InitS3Lib(lib.S3ConfigType{
		AccessKeyId:     config.Config.S3.AccessKeyId,
		SecretAccessKey: config.Config.S3.SecretAccessKey,
		DefaultRegion:   config.Config.S3.DefaultRegion,
		Bucket:          config.Config.S3.Bucket,
		CDNUrl:          config.Config.S3.CDNUrl,
	})

	app.Serve()
}
