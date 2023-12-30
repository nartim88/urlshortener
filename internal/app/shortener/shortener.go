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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nartim88/urlshortener/config"
	"github.com/nartim88/urlshortener/internal/controller/api/routers"
	"github.com/nartim88/urlshortener/internal/service"
	"github.com/nartim88/urlshortener/internal/storage"
	"github.com/nartim88/urlshortener/pkg/httpserver"
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
	server := httpserver.NewServer(a.Handler, a.Configs.RunAddr, 30*time.Second)

	logger.Log.Info().Msgf("running server on %s", a.Configs.RunAddr)

	idleConnsClosed := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		if err := server.Shutdown(); err != nil {
			logger.Log.Error().Stack().Err(err).Send()
		}
		close(idleConnsClosed)
	}()

	err := <-server.Notify()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Error().Stack().Err(err).Send()
	}

	store := a.Service.GetStore()

	switch s := store.(type) {
	case storage.StorageWithService:
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.Close(ctx); err != nil {
			logger.Log.Error().Stack().Err(err).Msg("error while closing db connection")
		}
		logger.Log.Info().Msg("db connection is closed")
	}

	<-idleConnsClosed

	close(a.Service.MarkAsDeletedCh())
	close(a.Service.MarkAsDeletedResultCh())

check:
	for {
		_, ok := <-a.Service.MarkAsDeletedCh()
		switch {
		case ok:
			logger.Log.Info().Msg("chan are closing")
			time.Sleep(1 * time.Second)
		default:
			break check
		}
	}

	logger.Log.Info().Msg("server is closed")
}

// initStorage инициализирует хранилище
func (a *Application) initStorage() (storage.Storage, error) {
	switch {
	case a.Configs.DatabaseDSN != "":
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		pool, err := pgxpool.New(ctx, a.Configs.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("error while connecting to db: %w", err)
		}

		s := storage.NewDBStorage(pool)

		if err = s.Bootstrap(ctx); err != nil {
			return nil, fmt.Errorf("error while creating tables in db: %w", err)
		}
		logger.Log.Info().Msg("db storage is initialized")
		return s, nil

	case a.Configs.FileStoragePath != "":
		s, err := storage.NewFileStorage(a.Configs.FileStoragePath)
		if err != nil {
			return nil, fmt.Errorf("error while creating file storage: %w", err)
		}
		logger.Log.Info().Msg("file storage is initialized")
		return s, nil

	default:
		s := storage.NewMemStorage()
		logger.Log.Info().Msg("memory storage is initialized")
		return s, nil
	}
}
