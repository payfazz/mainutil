package mainutil

import (
	"os"
)

// ErrorHandler used as parameter to github.com/payfazz/go-errors/errhandler.With.
// Print err using Eprint and exit the program with exit status 1.
func ErrorHandler(err error) {
	Eprint(err)
	os.Exit(1)
}
