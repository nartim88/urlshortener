package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

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
	return &Config{}
}

// ParseConfigs инициализация парсинга конфигов из окружения и флагов
func (conf *Config) ParseConfigs() error {
	logger := zerolog.New(os.Stdout)
	if err := conf.parseDotenv(); err != nil {
		return err
	}
	conf.parseFlags(logger)
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
func (conf *Config) parseFlags(logger zerolog.Logger) {
	flag.StringVar(&conf.RunAddr, "a", "", "address and port to run server")
	flag.StringVar(&conf.BaseURL, "b", "", "server address before shorten URL")
	flag.StringVar(&conf.LogLevel, "l", LogLevel, "log level")
	flag.StringVar(&conf.FileStoragePath, "f", "", "full file name for saving URLs")
	flag.StringVar(&conf.DatabaseDSN, "d", "", "database DSN")
	flag.Parse()

	logger.Info().
		Str("SERVER_ADDRESS", conf.RunAddr).
		Str("BASE_URL", conf.BaseURL).
		Str("LOG_LEVEL", conf.LogLevel).
		Str("FILE_STORAGE_PATH", conf.FileStoragePath).
		Str("DATABASE_DSN", conf.DatabaseDSN).
		Str("SECRET_KEY", func() string {
			if conf.SecretKey != "" {
				return "true"
			}
			return "false"
		}()).
		Msg("got flags:")
}

// parseDotenv загружает в окружение переменные из .env
func (conf *Config) parseDotenv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
