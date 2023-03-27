package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/parlaynu/studio1767-api/api/graph"
	"github.com/parlaynu/studio1767-api/internal/config"
	"github.com/parlaynu/studio1767-api/internal/middleware/auth"
	"github.com/parlaynu/studio1767-api/internal/service"
)

func New(cfg *config.Config, noauth bool) (http.Handler, error) {

	// create the top router with some useful middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// auth middleware handles login, logout, etc. workflows
	if noauth == false {
		mw, err := auth.NewAuthMiddleware(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed creating auth middleware: %w", err)
		}
		r.Use(mw)
	}

	// test endpoint
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	// the playground
	ph := playground.Handler("GraphQL playground", "/query")
	r.Handle("/play", ph)

	// the graphql server
	svc, err := service.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed creating service: %w", err)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(svc)}))
	r.Handle("/", srv)

	return r, nil
}
