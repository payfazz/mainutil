package mainutil

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// CancelOnInteruptSignal .
func (env *Env) CancelOnInteruptSignal(cancel context.CancelFunc) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, getInterruptSigs()...)
		sig := <-c
		signal.Stop(c)

		env.info().Print(fmt.Sprintf("Got signal %s\n", sig.String()))
		cancel()
	}()
}
