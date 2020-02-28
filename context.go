package mainutil

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// WaitForInterruptThenCancelContext .
func (env *Env) WaitForInterruptThenCancelContext(ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, getInterruptSigs()...)

	var sig os.Signal
	select {
	case sig = <-c:
	case <-ctx.Done():
	}

	signal.Stop(c)
	if sig != nil {
		env.InfoLogger().Print(fmt.Sprintf("Got signal %s\n", sig.String()))
	}

	cancel()
}
