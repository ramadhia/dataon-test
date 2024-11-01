package config

import (
	"encoding/json"
	"sync"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

type Config struct {
	App App
	Log Log
	DB  DB
}

type App struct {
	Name   string `default:"web-api" env:"SERVICE_NAME"`
	Scheme string `default:"http" env:"SERVICE_SCHEME"`
	Host   string `default:"0.0.0.0" env:"SERVICE_HOST"`
	Port   string `default:"14000" env:"SERVICE_PORT"`
}

type Log struct {
	Level  string `default:"INFO" env:"LOG_LEVEL"`
	Format string `default:"json" env:"LOG_FORMAT"`
	ByPath string `env:"LOG_BY_PATH"`
}

type DB struct {
	Client    string `default:"postgres" env:"DB_CLIENT"`
	Host      string `default:"0.0.0.0" env:"DB_HOST"`
	User      string `default:"root" env:"DB_USER"`
	Password  string `default:"true" required:"true" env:"DB_PASSWORD"`
	Port      string `default:"3306" env:"DB_PORT"`
	Database  string `default:"bosnetdb" env:"DB_DATABASE"`
	Migration struct {
		Autoload bool   `env:"DB_RUN_MIGRATION"`
		Path     string `default:"./database/migration" env:"DB_MIGRATION_PATH"`
	}
	Debug bool `default:"false" env:"DB_DEBUG"`
}

var config *Config
var configLock = &sync.Mutex{}

// Instance will return Config instance singleton
func Instance() Config {
	if config == nil {
		err := Load()
		if err != nil {
			panic(err)
		}
	}
	return *config
}

// Load will load configuration from env to default config instance
func Load() error {
	tmpConfig := Config{}
	err := configor.Load(&tmpConfig)
	if err != nil {
		return err
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &tmpConfig

	return nil
}

func (c Config) Json() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		logrus.WithError(err).Error("When marshalling Json")
	}
	return string(bytes)
}
