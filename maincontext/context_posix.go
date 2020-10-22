// +build linux darwin

package maincontext

import (
	"os"
	"syscall"
)

func getInterruptSigs() []os.Signal {
	return []os.Signal{syscall.SIGTERM, syscall.SIGINT}
}
