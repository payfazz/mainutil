// +build !linux,!darwin

package maincontext

import (
	"os"
)

func getInterruptSigs() []os.Signal {
	return []os.Signal{os.Interrupt}
}
