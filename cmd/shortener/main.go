package main

import (
	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/internal/pkg/routers"
)

func main() {
	app := shortener.New()
	app.Init()
	app.Run(routers.MainRouter())
}
