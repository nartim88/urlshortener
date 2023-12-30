package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nartim88/urlshortener/internal/controller/api/handlers"
	"github.com/nartim88/urlshortener/internal/controller/api/middleware"
	"github.com/nartim88/urlshortener/internal/service"
)

func MainRouter(srv service.Service) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.WithLogging)
	r.Use(middleware.GZipMiddleware)
	r.Use(middleware.AuthMiddleware(srv.GetConfigs()))

	r.Route("/", func(r chi.Router) {
		r.Mount("/", textRespRouter(srv))
	})

	r.Route("/api", func(r chi.Router) {
		r.Mount("/", apiRouter(srv))
	})

	r.Mount("/ping", dbPingRouter(srv))

	return r
}

func dbPingRouter(srv service.Service) http.Handler {
	r := chi.NewRouter()
	r.Get("/", handlers.DBPingHandle(srv))
	return r
}

func textRespRouter(srv service.Service) http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.IndexHandle(srv))
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetURLHandle(srv))
		})
	})
	return r
}

func apiRouter(srv service.Service) http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", handlers.GetShortURLHandle(srv))
			r.Route("/batch", func(r chi.Router) {
				r.Post("/", handlers.GetBatchShortURLsHandle(srv))
			})
		})
	})
	r.Mount("/user", userRouter(srv))
	return r
}

func userRouter(srv service.Service) http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Route("/urls", func(r chi.Router) {
			r.Get("/", handlers.UserURLsGet(srv))
			r.Delete("/", handlers.UserURLsDelete(srv))
		})
	})
	return r
}
