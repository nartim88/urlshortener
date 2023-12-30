package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewServer(h http.Handler, addr string, timeout time.Duration) *Server {
	httpServer := &http.Server{Handler: h, Addr: addr}
	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: timeout,
	}
	s.start()
	return s
}

func (s Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s Server) Notify() <-chan error {
	return s.notify
}

func (s Server) Shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	errCh := make(chan error)

	go func() {
		<-shutdownCtx.Done()
		if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
			err := fmt.Errorf("graceful shutdown timeout, forcing quit: %w", shutdownCtx.Err())
			errCh <- err
		}
	}()

	var err error
	select {
	case err = <-errCh:
	default:
		err = s.server.Shutdown(shutdownCtx)
	}

	return err
}
