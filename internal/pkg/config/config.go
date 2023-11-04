package config

import (
	"flag"

	"github.com/caarlos0/env"
	"github.com/nartim88/urlshortener/internal/pkg/logger"
)

type Config struct {
	RunAddr  string `env:"SERVER_ADDRESS"`
	BaseURL  string `env:"BASE_URL"`
	LogLevel string `env:"LOG_LEVEL"`
}

func New() *Config {
	cfg := Config{
		RunAddr:  "localhost:8080",
		BaseURL:  "http://localhost",
		LogLevel: "info",
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
		logger.Log.Info().Err(err).Send()
	}
}

// parseFlags парсинг флагов
func (conf *Config) parseFlags() {
	flag.StringVar(&conf.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "server address before shorten URL")
	flag.StringVar(&conf.LogLevel, "l", "info", "log level")

	flag.Parse()
}
