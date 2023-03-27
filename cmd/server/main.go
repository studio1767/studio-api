package main

import (
	"crypto/tls"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/parlaynu/studio1767-api/api"
	"github.com/parlaynu/studio1767-api/internal/config"
)

func init() {

}

func main() {
	// setup logging
	log.SetLevel(log.TraceLevel)
	formatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	}
	log.SetFormatter(formatter)

	// parse command line
	noauth := flag.Bool("noauth", false, "bypass authentication for the server")

	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Usage: %s [--noauth] <config-file>", filepath.Base(os.Args[0]))
	}
	cfgFile := flag.Arg(0)

	// load the configuration
	cfg, err := config.Load(cfgFile, *noauth)
	if err != nil {
		log.Fatal(err)
	}

	// create the service and api
	handler, err := api.New(cfg, *noauth)
	if err != nil {
		log.Fatal(err)
	}

	// run the server
	if strings.HasPrefix(cfg.Listener, "http://") {
		RunHTTP(cfg.Listener, handler)
	} else {
		RunHTTPS(cfg.Listener, handler, cfg.Https.CertFile, cfg.Https.KeyFile)
	}
}

func RunHTTP(listener string, handler http.Handler) {
	address := strings.TrimPrefix(listener, "http://")
	srv := &http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Infof("Server listening at %s\n", listener)
	log.Fatal(srv.ListenAndServe())
}

func RunHTTPS(listener string, handler http.Handler, certfile, keyfile string) {
	// create the TLS config
	// note: using tls 1.3 so all default ciphers etc. are secure
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	address := strings.TrimPrefix(listener, "https://")
	srv := &http.Server{
		Addr:      address,
		Handler:   handler,
		TLSConfig: tlsConfig,
		// chi middleware is handling timeouts
		// ReadTimeout:  15 * time.Second,
		// WriteTimeout: 15 * time.Second,
		// IdleTimeout:  60 * time.Second,
	}

	log.Infof("Server listening at %s\n", listener)
	log.Fatal(srv.ListenAndServeTLS(certfile, keyfile))
}
