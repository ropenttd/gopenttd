// openttd_adm_query is effectively a debugging tool that scans a given server over the admin port and dumps all available data.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	gopenttd "github.com/ropenttd/gopenttd/pkg/admin"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	serverHost  string
	serverPort  int
	serverPass  string
	prettyPrint bool
)

func init() {
	flag.StringVar(&serverHost, "target.host", "testserver.ttdredd.it", "Target host to connect to.")
	flag.IntVar(&serverPort, "target.port", 3977, "Target port (this should be the admin port)")
	flag.StringVar(&serverPass, "target.pass", "", "Target password")
	flag.BoolVar(&prettyPrint, "prettyprint", false, "Pretty print resulting JSON.")
	flag.Parse()
}

func main() {

	s, err := gopenttd.New(serverHost, serverPort, serverPass)
	if err != nil {
		log.Fatal(err)
	}
	s.LogLevel = gopenttd.LogInformational

	err = s.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// stupid delay to make sure things settle into state (? find a better way to do this)
	time.Sleep(3 * time.Second)

	state := s.State
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
