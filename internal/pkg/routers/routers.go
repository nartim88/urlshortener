package routers

import (
	"github.com/go-chi/chi/v5"

	"github.com/nartim88/urlshortener/internal/pkg/handlers"
	"github.com/nartim88/urlshortener/internal/pkg/middleware"
)

func MainRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.All...)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.IndexHandle)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetURLHandle)
		})
	})

	r.Mount("/api", apiRouter())

	r.Mount("/ping", dbPingRouter())

	return r
}

func apiRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/shorten", handlers.JSONGetShortURLHandle)
	})

	return r
}

func dbPingRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", handlers.DBPingHandle)
	return r
}
