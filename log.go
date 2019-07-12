package mainutil

import (
	"fmt"

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
