package auth

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/studio1767/studio-api/internal/config"
)

type Authenticator interface {
	Authenticate(ctx context.Context) (context.Context, error)
}

type GroupGetter interface {
	GroupsForUser(username string) (map[string]bool, error)
}

func NewAuthenticator(cfg *config.Config, gg GroupGetter) (Authenticator, error) {

	// wrap the group getter in a cache
	gg = NewCache(gg)

	// return the authenticator
	return &authenticator{
		gg: gg,
	}, nil
}

type authenticator struct {
	gg GroupGetter
}

type emailContextKey struct{}
type groupsContextKey struct{}
type getterContextKey struct{}

func (a *authenticator) Authenticate(ctx context.Context) (context.Context, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return ctx, status.New(codes.Unknown, "no peer information found").Err()
	}

	// verify we have auth info and it's of type tls
	if peer.AuthInfo == nil || peer.AuthInfo.AuthType() != "tls" {
		return ctx, status.New(codes.Unauthenticated, "no auth information found").Err()
	}

	// get the tls info to extract the common name and any groups
	tlsInfo := peer.AuthInfo.(credentials.TLSInfo)

	email := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
	ctx = context.WithValue(ctx, emailContextKey{}, email)

	groups := make(map[string]bool)
	for _, uri := range tlsInfo.State.VerifiedChains[0][0].URIs {
		if uri.Scheme == "group" {
			groups[uri.Opaque] = true
		}
	}

	if len(groups) > 0 {
		ctx = context.WithValue(ctx, groupsContextKey{}, groups)
	}

	// add the group getter to the context
	ctx = context.WithValue(ctx, getterContextKey{}, a.gg)

	return ctx, nil
}

func UnaryAuthnInterceptor(a Authenticator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, err := a.Authenticate(ctx)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func StreamAuthnInterceptor(a Authenticator) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx, err := a.Authenticate(stream.Context())
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
