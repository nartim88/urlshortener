package shortener

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nartim88/urlshortener/internal/config"
	"github.com/nartim88/urlshortener/internal/routers"
	"github.com/nartim88/urlshortener/internal/storage"
)

var CFG config.Config

type App struct {
	store   storage.Storage
	configs config.Config
}

func New() *App {
	store := storage.New()
	configs := config.New()
	return &App{
		store:   *store,
		configs: *configs,
	}
}

func (app *App) Init() *App {
	app.configs.Set()
	return app
}

func (app *App) Run() {
	fmt.Printf("Runnig server on %s", app.configs.RunAddr)

	err := http.ListenAndServe(app.configs.RunAddr, routers.MainRouter())

	if err != nil {
		log.Fatal(err)
	}
}
