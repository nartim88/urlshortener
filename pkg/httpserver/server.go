package httpserver

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewServer(h http.Handler, addr string) *Server {
	httpServer := &http.Server{Handler: h, Addr: addr}
	s := &Server{
		server: httpServer,
		notify: make(chan error, 1),
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
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}