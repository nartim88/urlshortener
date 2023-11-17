package shortener

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	logger.Log.Info().Msg("app configs:")
	logger.Log.Info().Str("SERVER_ADDRESS", a.Configs.RunAddr).Send()
	logger.Log.Info().Str("BASE_URL", a.Configs.BaseURL).Send()
	logger.Log.Info().Str("LOG_LEVEL", a.Configs.LogLevel).Send()
	logger.Log.Info().Str("FILE_STORAGE_PATH", a.Configs.FileStoragePath).Send()
	logger.Log.Info().Str("DATABASE_DSN", a.Configs.DatabaseDSN).Send()

	// инициализация хранилища
	store, err := a.initStorage()
	if err != nil {
		logger.Log.Error().Stack().Err(err).Send()
	}
	a.Store = store
}

// Run запуск сервера
func (a *Application) Run(h http.Handler) {
	var srv http.Server

	logger.Log.Info().Msgf("running server on %s", a.Configs.RunAddr)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Log.Error().Stack().Err(err).Send()
		}
		close(idleConnsClosed)
	}()

	srv.Addr = a.Configs.RunAddr
	srv.Handler = h

	err := srv.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		s, ok := a.Store.(storage.DBStorage)
		if ok {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := s.Close(ctx); err != nil {
				logger.Log.Error().Stack().Err(err).Msg("error while closing db connection")
			}
			logger.Log.Info().Msg("db connection is closed")
		}
	} else {
		logger.Log.Error().Stack().Err(err).Send()
	}

	<-idleConnsClosed
	logger.Log.Info().Msg("server is closed")
}

func (a *Application) initStorage() (storage.Storage, error) {
	switch {
	case a.Configs.DatabaseDSN != "":
		s, err := storage.NewDBStorage(a.Configs.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("error while creating db storage: %w", err)
		}
		return s, nil
	case a.Configs.FileStoragePath != "":
		s, err := storage.NewFileStorage(a.Configs.FileStoragePath)
		if err != nil {
			return nil, fmt.Errorf("error while creating file storage: %w", err)
		}
		return s, nil
	default:
		s := storage.NewMemStorage()
		return s, nil
	}
}
