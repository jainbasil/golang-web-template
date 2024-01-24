package config

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	Env       string         `mapstructure:"env"`
	Port      string         `mapstructure:"env"`
	Mongo     MongoDbConfig  `mapstructure:"mongodb"`
	RabbitMQ  RabbitMQConfig `mapstructure:"rabbitmq"`
	GcpConfig GcpConfig      `mapstructure:"gcp"`
}

type MongoDbConfig struct {
	Uri      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type RabbitMQConfig struct {
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Exchange   string `mapstructure:"exchange"`
	Queue      string `mapstructure:"queue"`
	FetchCount string `mapstructure:"fetch_count"` // only applicable for consumers
}

type GcpConfig struct {
	ProjectID    string `mapstructure:"project_id"`
	PubsubConfig struct {
		Topic        string `mapstructure:"topic"`
		Subscription string `mapstructure:"subscription"`
	} `mapstructure:"pubsub"`
}

func LoadConfig() *AppConfig {
	var appConfig AppConfig

	// reading env.yaml from .config directory
	// See resources/env-sample.yaml for sample file.
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("error loading env file from specified path: ", err)
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		log.Fatal("error unmarshalling appConfig")
	}

	// Use the following if you want to read one more config file
	//viper.SetConfigName("another-config")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath("./resources")
	//if err := viper.ReadInConfig(); err != nil {
	//	log.Fatal("error loading config from specified path: ", err)
	//}

	return &appConfig
}
