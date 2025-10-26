package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	Http        HttpConfig `yaml:"http" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
}

type HttpConfig struct {
	Port    int           `yaml:"port" env-default:"8000"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func MustLoad() Config {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
