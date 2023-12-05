package main

import (
	"github.com/nartim88/urlshortener/config"
	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/pkg/logger"
)

func main() {
	cfg := config.NewConfig()
	err := cfg.ParseConfigs()
	if err != nil {
		logger.Log.Error().Stack().Err(err).Send()
		panic(err)
	}
	app := shortener.NewApp()
	app.Init(cfg)
	app.Run()
}
