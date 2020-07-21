package mainutil

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

var (
	backgroundContext struct {
		once   sync.Once
		ctx    context.Context
		cancel context.CancelFunc
	}
)

// BackgroundContext return context that will be done when the program interrupted (SIGINT, SIGTERM)
func BackgroundContext() (context.Context, context.CancelFunc) {
	backgroundContext.once.Do(func() {
		backgroundContext.ctx, backgroundContext.cancel = context.WithCancel(context.Background())

		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, getInterruptSigs()...)

			select {
			case <-c:
			case <-backgroundContext.ctx.Done():
			}

			signal.Stop(c)
			backgroundContext.cancel()
		}()
	})
	return backgroundContext.ctx, backgroundContext.cancel
}
