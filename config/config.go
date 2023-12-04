package config

import (
	"flag"

	"github.com/caarlos0/env"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type Config struct {
	RunAddr         string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	SecretKey       string `env:"SECRET_KEY"`
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
func (conf *Config) ParseConfigs() error {
	if err := conf.parseDotenv(); err != nil {
		return err
	}
	conf.parseFlags()
	if err := conf.parseEnv(); err != nil {
		return err
	}
	return nil
}

// parseEnv парсит переменные окружения
func (conf *Config) parseEnv() error {
	err := env.Parse(conf)
	if err != nil {
		return err
	}
	return nil
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
func (conf *Config) parseDotenv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
