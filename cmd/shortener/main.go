package main

import (
	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/internal/pkg/routers"
)

func main() {
	shortener.App.Init()
	shortener.App.Run(routers.MainRouter())
}
