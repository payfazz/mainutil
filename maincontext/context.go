package maincontext

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

var (
	once sync.Once
	ctx  context.Context
)

// Background return context that will be done when the program got interrupted
func Background() context.Context {
	once.Do(func() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(context.Background())

		waitForInterrupt := func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, getInterruptSigs()...)
			<-c
			signal.Stop(c)
			cancel()
		}

		go waitForInterrupt()
	})
	return ctx
}
