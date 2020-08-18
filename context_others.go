// +build !linux,!darwin

package mainutil

import (
	"os"
)

func getInterruptSigs() []os.Signal {
	return []os.Signal{os.Interrupt}
}
