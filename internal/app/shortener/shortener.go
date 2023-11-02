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

var (
	App  Application
	Conf config.Config
	St   storage.Storage
)

type Application struct {
	Store   storage.Storage
	Configs config.Config
}

func NewApplication() {
	St = *storage.New()
	Conf = *config.New()
	App = Application{
		Store:   St,
		Configs: Conf,
	}
}

func (a *Application) Init() {
	NewApplication()

	if err := logger.Init(a.Configs.LogLevel); err != nil {
		logger.Log.Info().Stack().Err(err).Send()
	}

	a.Configs.Parse()
}

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
