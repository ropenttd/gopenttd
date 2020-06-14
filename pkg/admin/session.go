package admin

import (
	"time"
)

// VERSION of Gopenttd, follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.0.1"

// New creates a new OpenTTD Admin session and will automate some startup
// tasks if given enough information to do so.  Currently you can pass zero
// arguments and it will return an empty session.
func New(hostname string, port int, password string) (s *Session, err error) {

	// Create an empty Session interface.
	s = &Session{
		State:                  NewState(),
		StateEnabled:           true,
		ShouldReconnectOnError: true,
		UserAgent:              "gopenttd (https://github.com/ropenttd/gopenttd)",
		LastPong:               time.Now().UTC(),
		Hostname:               hostname,
		Port:                   port,
		Password:               password,
	}

	// You should now call Open() so that events will trigger.

	return
}
