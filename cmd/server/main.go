package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/studio1767/studio-api/internal/auth"
	"github.com/studio1767/studio-api/internal/config"
	"github.com/studio1767/studio-api/internal/db"
	"github.com/studio1767/studio-api/internal/ldapgroups"
	"github.com/studio1767/studio-api/internal/server"
)

func main() {
	// setup logging
	log.SetLevel(log.TraceLevel)
	formatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	}
	log.SetFormatter(formatter)

	// parse command line
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Usage: %s <config-file>", filepath.Base(os.Args[0]))
	}
	cfgFile := flag.Arg(0)

	// load the configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	// create the tls configs
	sTlsConfig, err := buildServerTlsConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	cTlsConfig, err := buildClientTlsConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// create the db client
	dbClient, err := db.NewClient(cfg, cTlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	// create the ldap client
	ldapClient, err := ldapgroups.NewClient(cfg, cTlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	// create the authenticator
	authenticator, err := auth.NewAuthenticator(cfg, ldapClient)
	if err != nil {
		log.Fatal(err)
	}

	// create the service
	srv, err := server.New(sTlsConfig, dbClient, authenticator)
	if err != nil {
		log.Fatal(err)
	}

	// create the listener
	listen := fmt.Sprintf("%s:%d", cfg.Service.ListenAddress, cfg.Service.ListenPort)
	fmt.Printf("listening at %s\n", listen)
	l, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatal(err)
	}

	// serve the api
	srv.Serve(l)
}

func buildServerTlsConfig(cfg *config.Config) (*tls.Config, error) {

	// create the TLS config
	tlsConfig := tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	// the servers certificate and key
	serverKeyPair, err := tls.LoadX509KeyPair(cfg.Service.CertFile, cfg.Service.KeyFile)
	if err != nil {
		return nil, err
	}

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

	return &tlsConfig, nil
}

func buildClientTlsConfig(cfg *config.Config) (*tls.Config, error) {
	// register the tls certificate with the mysql driver
	pem, err := os.ReadFile(cfg.Service.CaCertFile)
	if err != nil {
		return nil, err
	}
	certs := x509.NewCertPool()
	certs.AppendCertsFromPEM(pem)

	tlsConfig := tls.Config{
		RootCAs: certs,
	}

	return &tlsConfig, nil
}
