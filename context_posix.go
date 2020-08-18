// +build linux darwin

package mainutil

import (
	"os"
	"syscall"
)

func getInterruptSigs() []os.Signal {
	return []os.Signal{syscall.SIGTERM, syscall.SIGINT}
}
