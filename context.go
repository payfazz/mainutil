package mainutil

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/payfazz/stdlog"
)

// WaitForInterruptThenCancelContext .
func WaitForInterruptThenCancelContext(ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, getInterruptSigs()...)

	var sig os.Signal
	select {
	case sig = <-c:
	case <-ctx.Done():
	}

	signal.Stop(c)
	if sig != nil {
		stdlog.PrintOut(fmt.Sprintf("Got signal %s\n", sig.String()))
	}

	cancel()
}
