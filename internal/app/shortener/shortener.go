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

	"github.com/jackc/pgx/v5"
	"github.com/nartim88/urlshortener/config"
	"github.com/nartim88/urlshortener/internal/controller/api/routers"
	"github.com/nartim88/urlshortener/internal/service"
	"github.com/nartim88/urlshortener/internal/storage"
	"github.com/nartim88/urlshortener/pkg/logger"
)

type Application struct {
	Configs *config.Config
	Service service.Service
	Handler http.Handler
}

func NewApp() *Application {
	return &Application{}
}

// Init первичная инициализация приложения
func (a *Application) Init(cfg *config.Config) {

	// инициализация логгера
	if err := logger.Init(cfg.LogLevel); err != nil {
		logger.Log.Info().Stack().Err(err).Send()
	}

	logger.Log.Info().
		Str("SERVER_ADDRESS", cfg.RunAddr).
		Str("BASE_URL", cfg.BaseURL).
		Str("LOG_LEVEL", cfg.LogLevel).
		Str("FILE_STORAGE_PATH", cfg.FileStoragePath).
		Str("DATABASE_DSN", cfg.DatabaseDSN).
		Str("SECRET_KEY", func() string {
			if cfg.SecretKey != "" {
				return "true"
			}
			return "false"
		}()).
		Msg("app configs:")

	a.Configs = cfg

	// инициализация хранилища
	store, err := a.initStorage()
	if err != nil {
		logger.Log.Error().Stack().Err(err).Send()
	}
	a.Service = service.NewService(store, a.Configs)
	a.Handler = routers.MainRouter(a.Service)
}

// Run запуск сервера
func (a *Application) Run() {
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
	srv.Handler = a.Handler

	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Error().Stack().Err(err).Send()
	}

	store := a.Service.GetStore()
	s, ok := store.(storage.StorageWithService)
	if ok {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.Close(ctx); err != nil {
			logger.Log.Error().Stack().Err(err).Msg("error while closing db connection")
		}
		logger.Log.Info().Msg("db connection is closed")
	}

	<-idleConnsClosed
	logger.Log.Info().Msg("server is closed")
}

// initStorage инициализирует хранилище
func (a *Application) initStorage() (storage.Storage, error) {
	switch {
	case a.Configs.DatabaseDSN != "":
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		conn, err := pgx.Connect(ctx, a.Configs.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("error while connecting to db: %w", err)
		}

		s := storage.NewDBStorage(conn)

		if err = s.Bootstrap(ctx); err != nil {
			return nil, fmt.Errorf("error while creating tables in db: %w", err)
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
