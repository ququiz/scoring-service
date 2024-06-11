package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App
		HTTP
		LogConfig
		GRPC
		Mongodb
		Redis
		RabbitMQ
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		Port    string `env-required:"true" env:"HTTP_PORT"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	LogConfig struct {
		Level string `json:"level" yaml:"level" env:"LOG_LEVEL"`
		// Filename   string `json:"filename" yaml:"filename"`
		// MaxSize    int    `json:"maxsize" yaml:"maxsize"`
		MaxAge     int `json:"max_age" yaml:"max_age" env:"LOG_MAXAGE"`
		MaxBackups int `json:"max_backups" yaml:"max_backups" env:"LOG_MAXBACKUP"`
	}

	GRPC struct {
		URLGrpc        string `json:"urlGRPC"  env:"URL_GRPC"`
		QueryQueryGRPC string `json:"quizQueryGRPCURL"  env:"QUIZ_QUERY_GRPC_URL"`
		AuthGRPCClient string `json:"AuthGRPCClient" env:"AUTH_GRPC_CLIENT"`
	}

	Mongodb struct {
		MongoURL      string `json:"mongo_url"  env:"MONGO_URL"`
		Database      string `json:"mongo_db" env:"MONGO_DB"`
		MongoWriteURL string `json:"mongo_write_url"env:"MONGO_WRITE_URL"`
	}

	Redis struct {
		RedisAddr     string `json:"redis_addr" env:"REDIS_ADDR"`
		RedisPassword string `json:"redis_password" env:"REDIS_PASSWORD"`
	}

	RabbitMQ struct {
		RMQAddress string `json:"rabbitmqAddress"  env:"RABBITMQ_ADDRESS"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// err = cleanenv.ReadConfig(path+".env", cfg) // buat di doker , ../.env kalo debug (.env kalo docker)
	// err = cleanenv.ReadConfig(path+"/local.env", cfg) // local run
	if os.Getenv("APP_ENV") == "local" {
		err = cleanenv.ReadConfig(path+"/local.env", cfg)
	} else if os.Getenv("APP_ENV") == "k8s" {
		err = cleanenv.ReadConfig(path+"k8s.env", cfg)
	} else {
		err = cleanenv.ReadConfig(path+".env", cfg)
	}
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
