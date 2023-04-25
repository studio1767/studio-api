package server

import (
	"context"
	"fmt"

	api "github.com/parlaynu/studio1767-api/api/v1"
	"github.com/parlaynu/studio1767-api/internal/auth"
)

func (svr *studioServer) Ping(ctx context.Context, req *api.PingRequest) (*api.PingReply, error) {
	fmt.Println("Ping")

	email := auth.EmailFromContext(ctx)
	resp := api.PingReply{
		Message: fmt.Sprintf("ping %s %s!", req.Name, email),
	}

	return &resp, nil
}
