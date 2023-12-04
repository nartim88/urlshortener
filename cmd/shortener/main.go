package main

import (
	"github.com/nartim88/urlshortener/config"
	"github.com/nartim88/urlshortener/internal/app/shortener"
)

func main() {
	app := shortener.NewApp()
	cfg := config.NewConfig()
	_ = cfg.ParseConfigs()
	app.Init(cfg)
	app.Run()
}
