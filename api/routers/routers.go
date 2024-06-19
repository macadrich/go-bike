package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/macadrich/go-bike/api/handlers"
	"github.com/macadrich/go-bike/api/middleware"
	"github.com/macadrich/go-bike/config"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(handlers *handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition"
	))

	r.Group(func(r chi.Router) {
		r.Use(middleware.StaticTokenAuthorization(config.LoadAuthorization()))
		r.Route("/api/v1", func(r chi.Router) {
			r.Post("/indego-data-fetch-and-store-it-db", handlers.InsertStation)
			r.Get("/stations", handlers.QueryAllStation)
			r.Get("/stations/{kioskId}", handlers.QuerySpecificStation)
		})
	})

	r.Get("/healthcheck", handlers.HealthCheck)

	return r
}
