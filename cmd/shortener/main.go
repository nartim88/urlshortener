package main

import (
	"github.com/nartim88/urlshortener/config"
	"github.com/nartim88/urlshortener/internal/app/shortener"
)

func main() {
	cfg := config.NewConfig()
	err := cfg.ParseConfigs()
	if err != nil {
		panic(err)
	}
	app := shortener.NewApp()
	app.Init(cfg)
	app.Run()
}
