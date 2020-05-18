package main

import (
	"encoding/json"
	"flag"
	"fmt"
	gopenttd "github.com/ropenttd/gopenttd/pkg/admin"
	log "github.com/sirupsen/logrus"
)

var (
	serverHost   string
	serverPort   int
	serverPass   string
	logLevel     string
	prettyPrint  bool
	ignoreErrors bool
)

func init() {
	flag.StringVar(&serverHost, "target.host", "188.40.223.196", "Target host to connect to.")
	flag.IntVar(&serverPort, "target.port", 3977, "Target port (this should be the admin port)")
	flag.StringVar(&serverPass, "target.pass", "", "Target password")
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
	state, err := gopenttd.ScanServerAdm(serverHost, serverPort, serverPass)
	if err != nil && !ignoreErrors {
		log.Fatal(err)
	}
	var b []byte
	if prettyPrint {
		b, err = json.MarshalIndent(state, "", "    ")
	} else {
		b, err = json.Marshal(state)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
