package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddr string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

func New() *Config {
	cfg := Config{
		RunAddr: "localhost:8080",
		BaseURL: "http://localhost",
	}
	return &cfg
}

// Parse инициализация конфигов приложения
func (conf *Config) Parse() {
	conf.parseFlags()
	conf.parseEnv()
}

// parseEnv парсинг переменных окружения
func (conf *Config) parseEnv() {
	err := env.Parse(conf)
	if err != nil {
		log.Fatal(err)
	}
}

// parseFlags парсинг флагов
func (conf *Config) parseFlags() {
	flag.StringVar(&conf.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "server address before shorten URL")
	flag.Parse()
}
