package mainutil

import (
	"os"

	"github.com/payfazz/go-errors"
)

// ErrorHandler used as parameter to github.com/payfazz/go-errors/errhandler.With.
// Print err using Eprint and exit the program with exit status 1.
func ErrorHandler(err error) {
	errors.PrintTo(Err, errors.Wrap(err))
	os.Exit(1)
}
