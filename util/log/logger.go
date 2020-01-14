package log

import (
	stdLog "log"
)

// Enabled allow log output to console
var Enabled bool

//Debugged enable output debug log
var Debugged bool

// Printf delegates to log.Printf if Enabled == true
func Printf(msg string, args ...interface{}) {
	if Enabled {
		stdLog.Printf(msg, args)
	}
}
