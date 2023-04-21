package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/parlaynu/studio1767-api/internal/config"
	"github.com/parlaynu/studio1767-api/internal/server"
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

	// create the service
	srv, err := server.New(cfg)
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
