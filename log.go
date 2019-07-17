package mainutil

import (
	"fmt"
	"time"

	"github.com/payfazz/go-errors"
	"github.com/payfazz/stdlog"
)

var (
	// Out from stdlog.
	Out = stdlog.Out

	// Err from stdlog.
	Err = stdlog.Err
)

// Eprint print errors to stderr, comply with 12factor.net
func Eprint(err error) {
	Err.Print(errors.Format(err))
}

// EprintTime print errors to stderr, prefix it with UTC time, comply with 12factor.net
func EprintTime(err error) {
	Err.Print(time.Now().UTC(), ": ", errors.Format(err))
}

// Iprintf print information to stdout, comply with 12factor.net
func Iprintf(f string, v ...interface{}) {
	Out.Print(fmt.Sprintf(f, v...))
}

// IprintfTime print information to stdout, prefix it with UTC time, comply with 12factor.net
func IprintfTime(f string, v ...interface{}) {
	Out.Print(time.Now().UTC(), ": ", fmt.Sprintf(f, v...))
}
