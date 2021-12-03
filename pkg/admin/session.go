package admin

import (
	"github.com/ropenttd/gopenttd/pkg/admin/enum"
	"time"
)

// VERSION of Gopenttd, follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.1.0"

// New creates a new OpenTTD Admin session and will automate some startup
// tasks if given enough information to do so.  Currently you can pass zero
// arguments and it will return an empty session.
func New(hostname string, port int, password string) (s *Session, err error) {

	// Create an empty Session interface.
	s = &Session{
		State:                  NewState(),
		StateEnabled:           true,
		ShouldReconnectOnError: true,
		UpdateFrequencies:      map[enum.UpdateType]enum.UpdateFrequency{},
		UserAgent:              "gopenttd (https://github.com/ropenttd/gopenttd)",
		LastPong:               time.Now().UTC(),
		Hostname:               hostname,
		Port:                   port,
		Password:               password,
		rconQueue:              make(chan *rconRequest),
		rconChan:               make(chan *rconResp),
	}

	// You should now call Open() so that events will trigger.

	return
}
