package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/parlaynu/studio1767-api/internal/service"
)

type Resolver struct {
	svc service.Service
}

func NewResolver(svc service.Service) *Resolver {
	return &Resolver{
		svc: svc,
	}
}
