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

	printErr(rec)

	os.Exit(1)
}

func printErr(data interface{}) {
	if data == nil {
		return
	}
	if err, ok := data.(error); ok {
		errors.PrintTo(stdlog.Err(), errors.Wrap(err))
	} else {
		errors.PrintTo(stdlog.Err(), errors.Errorf("unknown error: %#v", data))
	}
}
