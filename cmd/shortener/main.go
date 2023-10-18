package main

import (
	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/internal/routers"
)

func main() {
	shortener.New()
	shortener.App.Init()
	shortener.App.Run(routers.MainRouter())
}
