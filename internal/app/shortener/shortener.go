package shortener

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

func (a *Application) Run(h http.Handler) {
	var srv http.Server

	log.Printf("Runnig server on %s", a.Configs.RunAddr)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	srv.Addr = a.Configs.RunAddr
	srv.Handler = h

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Print("Server closed.")
}
