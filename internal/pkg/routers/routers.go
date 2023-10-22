package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/nartim88/urlshortener/internal/pkg/handlers"
)

func MainRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.IndexHandle)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetURLHandle)
		})
	})

	return r
}
