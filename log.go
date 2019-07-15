package mainutil

import (
	"fmt"
	"time"

	"github.com/payfazz/go-errors"
	"github.com/payfazz/stdlog"
)

// Eprint print errors to stderr, comply with 12factor.net
func Eprint(err error) {
	stdlog.E(errors.Format(err))
}

// Iprintf print information to stdout, comply with 12factor.net
func Iprintf(f string, v ...interface{}) {
	stdlog.O(fmt.Sprintf(f, v...))
}

// EprintTime print errors to stderr, prefix it with UTC time, comply with 12factor.net
func EprintTime(err error) {
	stdlog.E(time.Now().UTC(), ": ", errors.Format(err))
}

// IprintfTime print information to stdout, prefix it with UTC time, comply with 12factor.net
func IprintfTime(f string, v ...interface{}) {
	stdlog.O(time.Now().UTC(), ": ", fmt.Sprintf(f, v...))
}
