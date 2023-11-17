package shortener

import (
	"context"
	"errors"
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

// Init первичная инициализация приложения
func (a *Application) Init() {

	// инициализация конфигов
	a.Configs = *config.NewConfig()
	a.Configs.ParseConfigs()

	// инициализация логгера
	if err := logger.Init(a.Configs.LogLevel); err != nil {
		logger.Log.Info().Stack().Err(err).Send()
	}

	logger.Log.Info().Msg("App configs:")
	logger.Log.Info().Str("SERVER_ADDRESS", a.Configs.RunAddr).Send()
	logger.Log.Info().Str("BASE_URL", a.Configs.BaseURL).Send()
	logger.Log.Info().Str("LOG_LEVEL", a.Configs.LogLevel).Send()
	logger.Log.Info().Str("FILE_STORAGE_PATH", a.Configs.FileStoragePath).Send()
	logger.Log.Info().Str("DATABASE_DSN", a.Configs.DatabaseDSN).Send()

	// инициализация хранилища
	store, err := a.initStorage()
	if err != nil {
		logger.Log.Info().Err(err).Send()
	}
	a.Store = store
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

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Info().Stack().Err(err).Send()
	}

	<-idleConnsClosed
	logger.Log.Info().Msg("Server closed.")
}

func (a *Application) initStorage() (storage.Storage, error) {
	switch {
	case a.Configs.DatabaseDSN != "":
		s, err := storage.NewDBStorage(a.Configs.DatabaseDSN)
		return s, err
	case a.Configs.FileStoragePath != "":
		s, err := storage.NewFileStorage(a.Configs.FileStoragePath)
		return s, err
	default:
		s := storage.NewMemStorage()
		return s, nil
	}
}
