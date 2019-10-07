package mainutil

import (
	"fmt"
	"log"

	"github.com/payfazz/go-errors"
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

// Eprint print errors to stderr
func Eprint(err error) {
	Err.Print(errors.Format(errors.Wrap(err)))
}

// Iprintf print information to stdout
func Iprintf(f string, v ...interface{}) {
	Out.Print(fmt.Sprintf(f, v...))
}
