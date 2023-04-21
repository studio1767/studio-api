package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	api "github.com/parlaynu/studio1767-api/api/v1"
	"github.com/parlaynu/studio1767-api/internal/config"
)

func New(cfg *config.Config, opts ...grpc.ServerOption) (*grpc.Server, error) {

	// create the grpc server
	creds, err := buildTlsCredentials(cfg)
	if err != nil {
		return nil, err
	}

	opts = append(opts,
		creds,
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_auth.StreamServerInterceptor(authenticate),
			),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_auth.UnaryServerInterceptor(authenticate),
			),
		),
	)
	gsrv := grpc.NewServer(opts...)

	// create the studio server
	srv, err := newServer(cfg)
	if err != nil {
		return nil, err
	}

	// and register the studio server
	api.RegisterStudioServer(gsrv, srv)

	return gsrv, nil
}

type studioServer struct {
	api.UnimplementedStudioServer
	db *sql.DB
}

func newServer(cfg *config.Config) (*studioServer, error) {

	// register the tls certificate with the mysql driver
	pem, err := os.ReadFile(cfg.Service.CaCertFile)
	if err != nil {
		return nil, err
	}
	certs := x509.NewCertPool()
	certs.AppendCertsFromPEM(pem)
	tlsCfg := &tls.Config{
		RootCAs: certs,
	}
	mysql.RegisterTLSConfig("maria", tlsCfg)

	// create the db config
	dbCfg := mysql.NewConfig()
	dbCfg.Net = "tcp"
	dbCfg.Addr = fmt.Sprintf("%s:%d", cfg.Db.Server, cfg.Db.Port)
	dbCfg.DBName = cfg.Db.DbName
	dbCfg.User = cfg.Db.UserName
	dbCfg.Passwd = cfg.Db.Password
	dbCfg.TLSConfig = "maria"

	fmt.Println(dbCfg.FormatDSN())

	// connect to the database
	db, err := sql.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	svc := &studioServer{
		db: db,
	}

	return svc, nil
}

func buildTlsCredentials(cfg *config.Config) (grpc.ServerOption, error) {

	// create the TLS config
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	// the servers certificate and key
	serverKeyPair, err := tls.LoadX509KeyPair(cfg.Service.CertFile, cfg.Service.KeyFile)

	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0] = serverKeyPair
	if err != nil {
		return nil, err
	}

	// the ca certificate to authenticate the client (mTLS)
	caCert, err := os.ReadFile(cfg.Service.CaCertFile)
	if err != nil {
		return nil, err
	}
	ca := x509.NewCertPool()
	ok := ca.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf("failed to parse ca certificate")
	}

	tlsConfig.ClientCAs = ca
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert

	// create the grpc credentials
	creds := grpc.Creds(credentials.NewTLS(tlsConfig))

	return creds, nil
}

type emailContextKey struct{}
type groupsContextKey struct{}

func authenticate(ctx context.Context) (context.Context, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return ctx, status.New(
			codes.Unknown,
			"no peer information found",
		).Err()
	}

	// verify we have auth info and it's of type tls
	if peer.AuthInfo == nil || peer.AuthInfo.AuthType() != "tls" {
		return context.WithValue(ctx, emailContextKey{}, ""), nil
	}

	// get the tls info to extract the common name and any groups
	tlsInfo := peer.AuthInfo.(credentials.TLSInfo)

	email := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
	ctx = context.WithValue(ctx, emailContextKey{}, email)

	var groups []string
	for _, uri := range tlsInfo.State.VerifiedChains[0][0].URIs {
		if uri.Scheme == "group" {
			groups = append(groups, uri.Opaque)
		}
	}

	if len(groups) > 0 {
		ctx = context.WithValue(ctx, groupsContextKey{}, groups)
	}

	return ctx, nil
}
