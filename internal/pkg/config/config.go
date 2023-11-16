package config

import (
	"flag"
	"fmt"

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
	DBHost          string `env:"DB_HOST"`
	DBPort          string `env:"DB_PORT"`
	DBUser          string `env:"DB_USER"`
	DBPassword      string `env:"DB_PASSWORD"`
	DBName          string `env:"DB_NAME"`
}

// NewConfig инициализирует Config с дефолтными значениями
func NewConfig() *Config {
	cfg := Config{
		RunAddr:         "localhost:8080",
		BaseURL:         "http://localhost",
		LogLevel:        "info",
		FileStoragePath: "/tmp/short-url-db.json",
		DatabaseDSN:     "",
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
	var databaseDSN = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName)

	flag.StringVar(&conf.RunAddr, "a", RunAddr, "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", BaseURL, "server address before shorten URL")
	flag.StringVar(&conf.LogLevel, "l", LogLevel, "log level")
	flag.StringVar(&conf.FileStoragePath, "f", FileStoragePath, "full file name for saving URLs")
	flag.StringVar(&conf.DatabaseDSN, "d", databaseDSN, "database DSN")

	flag.Parse()
}

func (conf *Config) parseDotenv() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Info().Err(err).Send()
	}
}
