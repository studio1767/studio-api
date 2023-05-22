package server

import (
	"crypto/tls"
	"database/sql"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	api "github.com/studio1767/studio-api/api/v1"
	"github.com/studio1767/studio-api/internal/auth"
)

func New(sTlsConfig *tls.Config, dbClient *sql.DB, authn auth.Authenticator, opts ...grpc.ServerOption) (*grpc.Server, error) {

	// create the grpc server with TLS credentials and interceptors
	opts = append(opts,
		grpc.Creds(credentials.NewTLS(sTlsConfig)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			auth.StreamAuthnInterceptor(authn),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			auth.UnaryAuthnInterceptor(authn),
		)),
	)
	gsrv := grpc.NewServer(opts...)

	// create the studio server
	srv, err := newServer(dbClient)
	if err != nil {
		return nil, err
	}

	// and register the studio server
	api.RegisterStudioServer(gsrv, srv)

	return gsrv, nil
}

type studioServer struct {
	api.UnimplementedStudioServer
	dbClient *sql.DB
}

func newServer(dbClient *sql.DB) (*studioServer, error) {

	svc := &studioServer{
		dbClient: dbClient,
	}

	return svc, nil
}
