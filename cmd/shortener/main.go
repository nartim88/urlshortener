package main

import (
	"github.com/nartim88/urlshortener/config"
	"github.com/nartim88/urlshortener/internal/app/shortener"
)

func main() {
	cfg := config.NewConfig()
	cfg.ParseConfigs()
	app := shortener.NewApp()
	app.Init(cfg)
	app.Run()
}
