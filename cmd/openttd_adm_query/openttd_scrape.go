package main

import (
	"flag"
	log "github.com/sirupsen/logrus"

	"github.com/ropenttd/gopenttd/pkg/gopenttd"
)

var (
	serverHost   string
	serverPort   int
	logLevel     string
	prettyPrint  bool
	ignoreErrors bool
)

func init() {
	flag.StringVar(&serverHost, "target.host", "188.40.223.196", "Show help menu")
	flag.IntVar(&serverPort, "target.port", 3977, "Show version information")
	flag.StringVar(&logLevel, "loglevel", "debug", "Set log level.")
	flag.BoolVar(&prettyPrint, "prettyprint", false, "Pretty print resulting JSON.")
	flag.BoolVar(&ignoreErrors, "ignore-errors", false, "Don't exit on connection errors and always output JSON.")
	flag.Parse()
}

func main() {
	parsedLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(parsedLevel)
	err = gopenttd.ScanServerAdm(serverHost, serverPort)
	if err != nil && !ignoreErrors {
		log.Fatal(err)
	}
}
