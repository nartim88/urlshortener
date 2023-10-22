package main

import (
	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/internal/pkg/logger"
	"github.com/nartim88/urlshortener/internal/pkg/routers"
)

func main() {
	if err := logger.Init(shortener.App.Configs.LogLevel); err != nil {
		logger.Log.Info().Err(err)
	}
	shortener.New()
	shortener.App.Init()
	shortener.App.Run(routers.MainRouter())
}
