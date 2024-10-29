package config

import (
	"log"

	"github.com/spf13/viper"
)

var Config ConfigType

type ConfigType struct {
	Env    string `mapstructure:"env"`
	Port   int    `mapstructure:"port"`
	Debug  bool   `mapstructure:"debug"`
	Secret string `mapstructure:"secret"`
	SSO    struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"sso"`
	Host      string `mapstructure:"host"`
	FrontHost string `mapstructure:"fronthost"`
	Database  struct {
		URL        string `mapstructure:"url"`
		SqlDir     string `mapstructure:"sqldir"`
		Migrations string `mapstructure:"migrations"`
	} `mapstructure:"database"`
	Sendgrid struct {
		Disabled  bool              `mapstructure:"disabled"`
		URL       string            `mapstructure:"url"`
		ApiKey    string            `mapstructure:"api_key"`
		Templates map[string]string `mapstructure:"templates"`
	} `mapstructure:"sendgrid"`
	Wellet struct {
		Agent       string `mapstructure:"agent"`
		AgentApiKey string `mapstructure:"agent_api_key"`
		Connect     string `mapstructure:"connect"`
	} `mapstructure:"wallet"`
	S3 struct {
		AccessKeyId     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		DefaultRegion   string `mapstructure:"default_region"`
		Bucket          string `mapstructure:"bucket"`
		CDNUrl          string `mapstructure:"cdn_url"`
	} `mapstructure:"s3"`
	Cors struct {
		Origins []string `mapstructure:"origins"`
	} `mapstructure:"cors"`
	Nats struct {
		Url   string `mapstructure:"url"`
		Token string `mapstructure:"token"`
	} `mapstructure:"nats"`
	Logger struct {
		Discord map[string]string `mapstructure:"discord"`
	} `mapstructure:"logger"`
	Admin struct {
		AccessToken string `mapstructure:"access_token"`
	} `mapstructure:"admin"`
}

func Init(configPath string) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found: %s", err)
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatal(err)
	}

	log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
}
