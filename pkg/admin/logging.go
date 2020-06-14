package admin

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

const (

	// LogError level is used for critical errors that could lead to data loss
	// or panic that would not be returned to a calling function.
	LogError int = iota

	// LogWarning level is used for very abnormal events and errors that are
	// also returned to a calling function.
	LogWarning

	// LogInformational level is used for normal non-error activity
	LogInformational

	// LogDebug level is for very detailed non-error activity.  This is
	// very spammy and will impact performance.
	LogDebug
)

// Logger can be used to replace the standard logging for gopenttd
var Logger func(msgL, caller int, format string, a ...interface{})

// msglog provides package wide logging consistancy for gopenttd
// the format, a...  portion this command follows that of fmt.Printf
//   msgL   : LogLevel of the message
//   caller : 1 + the number of callers away from the message source
//   format : Printf style message format
//   a ...  : comma separated list of values to pass
func msglog(msgL, caller int, format string, a ...interface{}) {

	if Logger != nil {
		Logger(msgL, caller, format, a...)
	} else {

		pc, file, line, _ := runtime.Caller(caller)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := fmt.Sprintf(format, a...)

		log.Printf("[DG%d] %s:%d:%s() %s\n", msgL, file, line, name, msg)
	}
}

// helper function that wraps msglog for the Session struct
// This adds a check to insure the message is only logged
// if the session log level is equal or higher than the
// message log level
func (s *Session) log(msgL int, format string, a ...interface{}) {

	if msgL > s.LogLevel {
		return
	}

	msglog(msgL, 2, format, a...)
}
