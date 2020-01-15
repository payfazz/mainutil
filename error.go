package mainutil

import (
	"os"

	"github.com/payfazz/go-errors"
)

// ErrorHandler used as parameter to github.com/payfazz/go-errors/errhandler.With.
//
// print to Err and exit with exit code 1
func ErrorHandler(err error) {
	errors.PrintTo(Err, errors.Wrap(err))
	os.Exit(1)
}

// ErrHandlerPrintOnly used as parameter to github.com/payfazz/go-errors/errhandler.With.
//
// print to Err and exit with exit code 1
func ErrHandlerPrintOnly(err error) {
	errors.PrintTo(Err, errors.Wrap(err))
}
