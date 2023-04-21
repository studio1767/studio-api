package server

import (
	"context"
	"fmt"

	api "github.com/parlaynu/studio1767-api/api/v1"
)

func (svr *studioServer) Hello(ctx context.Context, req *api.HelloRequest) (*api.HelloReply, error) {

	email := ctx.Value(emailContextKey{}).(string)
	fmt.Printf("Hello %s\n", email)

	resp := api.HelloReply{
		Message: fmt.Sprintf("Hello %s %s!", req.Name, email),
	}

	return &resp, nil
}
