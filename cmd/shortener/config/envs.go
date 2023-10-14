package config

import (
	"log"

	"github.com/caarlos0/env"
)

// parseEnv парсит переменные окружения
func parseEnv() {
	err := env.Parse(&CFG)
	if err != nil {
		log.Fatal(err)
	}
}
