package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/nartim88/urlshortener/internal/pkg/handlers"
	"github.com/nartim88/urlshortener/internal/pkg/middleware"
)

func MainRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.WithLogging)
	r.Use(middleware.GZipMiddleware)
	r.Use(middleware.AuthMiddleware)

	r.Route("/", func(r chi.Router) {
		r.Mount("/", textRespRouter())
	})

	r.Route("/api", func(r chi.Router) {
		r.Mount("/", apiRouter())
	})

	r.Mount("/ping", dbPingRouter())

	return r
}

func dbPingRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", handlers.DBPingHandle)
	return r
}

func textRespRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.IndexHandle)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetURLHandle)
		})
	})
	return r
}

func apiRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", handlers.GetShortURLHandle)
			r.Route("/batch", func(r chi.Router) {
				r.Post("/", handlers.GetBatchShortURLsHandle)
			})
		})
	})
	r.Mount("/user", userRouter())
	return r
}

func userRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Route("/urls", func(r chi.Router) {
			r.Get("/", handlers.GetAllUserURLs)
		})
	})
	return r
}
