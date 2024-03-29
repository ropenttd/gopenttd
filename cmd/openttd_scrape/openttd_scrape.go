package main

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"

	gopenttd "github.com/ropenttd/gopenttd/pkg/query"
)

var (
	serverHost   string
	serverPort   int
	logLevel     string
	prettyPrint  bool
	ignoreErrors bool
)

func init() {
	flag.StringVar(&serverHost, "target.host", "localhost", "Target hostname or IP address.")
	flag.IntVar(&serverPort, "target.port", 3979, "Target port (this should be the game connection port).")
	flag.StringVar(&logLevel, "loglevel", "warn", "Set log level.")
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
	serverData, err := gopenttd.ScanServer(serverHost, serverPort)
	if err != nil && !ignoreErrors {
		log.Fatal(err)
	}
	var b []byte
	if prettyPrint {
		b, err = json.MarshalIndent(serverData, "", "    ")
	} else {
		b, err = json.Marshal(serverData)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
