package shortener

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nartim88/urlshortener/internal/pkg/config"
	"github.com/nartim88/urlshortener/internal/pkg/logger"
	"github.com/nartim88/urlshortener/internal/pkg/storage"
)

type Application struct {
	Store   storage.Storage
	Configs config.Config
}

var App Application

func newApplication() {
	Conf := config.NewConfig()
	Store, err := storage.NewFileStorage(Conf.FileStoragePath)
	if err != nil {
		logger.Log.Info().Err(err).Send()
	}
	App = Application{
		Store:   Store,
		Configs: *Conf,
	}
}

// Init первичная инициализация приложения
func (a *Application) Init() {
	newApplication()

	if err := logger.Init(a.Configs.LogLevel); err != nil {
		logger.Log.Info().Stack().Err(err).Send()
	}

	a.Configs.Parse()
}

// Run запуск сервера
func (a *Application) Run(h http.Handler) {
	var srv http.Server

	logger.Log.Info().Msgf("Running server on %s", a.Configs.RunAddr)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Log.Info().Stack().Err(err).Send()
		}
		close(idleConnsClosed)
	}()

	srv.Addr = a.Configs.RunAddr
	srv.Handler = h

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log.Info().Stack().Err(err).Send()
	}

	<-idleConnsClosed
	logger.Log.Info().Msg("Server closed.")
}
