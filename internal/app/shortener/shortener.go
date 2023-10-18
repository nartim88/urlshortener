package shortener

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/nartim88/urlshortener/internal/config"
	"github.com/nartim88/urlshortener/internal/storage"
)

var (
	App  Application
	Conf config.Config
	St   storage.Storage
)

type Application struct {
	Store   storage.Storage
	Configs config.Config
}

func New() {
	St = *storage.New()
	Conf = *config.New()
	App = Application{
		Store:   St,
		Configs: Conf,
	}
}

func (a *Application) Init() *Application {
	a.Configs.Parse()
	return a
}

func (a *Application) Run(router chi.Router) {
	fmt.Printf("Runnig server on %s", a.Configs.RunAddr)

	err := http.ListenAndServe(a.Configs.RunAddr, router)

	if err != nil {
		log.Fatal(err)
	}
}
