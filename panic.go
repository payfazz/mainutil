package mainutil

import (
	"os"

	"github.com/payfazz/go-errors"
	"github.com/payfazz/stdlog"
)

// ExitOnPanic .
func ExitOnPanic() {
	rec := recover()
	if rec == nil {
		return
	}

	if err, ok := rec.(error); ok {
		errors.PrintTo(stdlog.Err(), errors.Wrap(err))
	} else {
		errors.PrintTo(stdlog.Err(), errors.Errorf("unknown error: %#v", rec))
	}

	os.Exit(1)
}
