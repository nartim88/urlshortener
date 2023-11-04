package config

import (
	"flag"

	"github.com/caarlos0/env"
	"github.com/nartim88/urlshortener/internal/pkg/logger"
)

type Config struct {
	RunAddr         string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

// New инициализирует Config с дефолтными значениями
func New() *Config {
	cfg := Config{
		RunAddr:         "localhost:8080",
		BaseURL:         "http://localhost",
		LogLevel:        "info",
		FileStoragePath: "/tmp/short-url-db.json",
	}
	return &cfg
}

// Parse инициализация парсинга конфигов из окружения и флагов
func (conf *Config) Parse() {
	conf.parseFlags()
	conf.parseEnv()
}

// parseEnv парсит переменные окружения
func (conf *Config) parseEnv() {
	err := env.Parse(conf)
	if err != nil {
		logger.Log.Info().Err(err).Send()
	}
}

// parseFlags парсит флаги командной строки
func (conf *Config) parseFlags() {
	flag.StringVar(&conf.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "server address before shorten URL")
	flag.StringVar(&conf.LogLevel, "l", "info", "log level")
	flag.StringVar(&conf.FileStoragePath, "f", "/tmp/short-url-db.json", "full file name for saving URLs")

	flag.Parse()
}
