package mainutil

import (
	"fmt"

	"github.com/payfazz/go-errors"
	"github.com/payfazz/stdlog"
)

// Eprint .
func Eprint(err error) {
	stdlog.E(errors.Format(err))
}

// Iprintf .
func Iprintf(f string, v ...interface{}) {
	stdlog.O(fmt.Sprintf(f, v...))
}
