package config

import (
	"flag"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"

	"github.com/nartim88/urlshortener/internal/pkg/logger"
)

type Config struct {
	RunAddr         string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

// NewConfig инициализирует Config с дефолтными значениями
func NewConfig() *Config {
	cfg := Config{
		RunAddr:  "localhost:8080",
		BaseURL:  "http://localhost",
		LogLevel: "info",
	}
	return &cfg
}

// ParseConfigs инициализация парсинга конфигов из окружения и флагов
func (conf *Config) ParseConfigs() {
	conf.parseDotenv()
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
	flag.StringVar(&conf.RunAddr, "a", RunAddr, "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", BaseURL, "server address before shorten URL")
	flag.StringVar(&conf.LogLevel, "l", LogLevel, "log level")
	flag.StringVar(&conf.FileStoragePath, "f", "", "full file name for saving URLs")
	flag.StringVar(&conf.DatabaseDSN, "d", "", "database DSN")

	flag.Parse()
}

// parseDotenv загружает в окружение переменные из .env
func (conf *Config) parseDotenv() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Info().Err(err).Send()
	}
}
