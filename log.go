package mainutil

import (
	"fmt"
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
	Err = log.New(stdlog.Err, "", log.LstdFlags|log.LUTC|log.Lshortfile)
}

// Eprintf print errors to stderr
func Eprintf(f string, v ...interface{}) {
	Err.Print(fmt.Sprintf(f, v...))
}

// Iprintf print information to stdout
func Iprintf(f string, v ...interface{}) {
	Out.Print(fmt.Sprintf(f, v...))
}
