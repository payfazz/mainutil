package mainutil

import (
	"log"

	"github.com/payfazz/stdlog"
)

var (
	// Out to log to stdout
	Out *log.Logger

	// Err to log to stderr
	Err *log.Logger
)

func init() {
	Out = log.New(stdlog.Out, "", log.LstdFlags|log.LUTC)
	Err = log.New(stdlog.Err, "", log.LstdFlags|log.LUTC)
}

// Eprintf print errors to stderr
func Eprintf(f string, v ...interface{}) {
	Err.Printf(f, v...)
}

// Iprintf print information to stdout
func Iprintf(f string, v ...interface{}) {
	Out.Printf(f, v...)
}
