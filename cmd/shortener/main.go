package main

import (
	"github.com/nartim88/urlshortener/internal/app/shortener"
)

func main() {
	app := shortener.New()
	app.Init()
	app.Run()
}
