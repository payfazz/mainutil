package mainutil

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/payfazz/go-errors"
	"github.com/payfazz/go-errors/errhandler"
	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
	"github.com/payfazz/go-middleware/common/logger"
	"github.com/payfazz/go-middleware/common/paniclogger"
)

// RunHTTPServer run *http.Server,
// when SIGTERM or SIGINT is recieved graceful shutdown the server.
//
// github.com/payfazz/go-errors/errhandler.With must be already defered.
func RunHTTPServer(server *http.Server) {
	serverErrCh := make(chan error, 1)
	go func() {
		defer close(serverErrCh)
		if server.TLSConfig == nil {
			Iprintf("Server listen on \"%s\"\n", server.Addr)
			serverErrCh <- errors.Wrap(server.ListenAndServe())
		} else {
			Iprintf("Server listen on TLS \"%s\"\n", server.Addr)
			serverErrCh <- errors.Wrap(server.ListenAndServeTLS("", ""))
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	signal.Notify(signalChan, signals...)

	select {
	case err := <-serverErrCh:
		signal.Reset(signals...)
		errhandler.Fail(errors.Wrap(err))
	case sig := <-signalChan:
		signal.Reset(signals...)
		waitFor := (1 * time.Minute) + (30 * time.Second)
		Iprintf(
			"Got '%s' signal, Shutting down the server (Waiting for graceful shutdown: %s)\n",
			sig.String(), waitFor.String(),
		)
		ctx, cancel := context.WithTimeout(context.Background(), waitFor)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			errhandler.Fail(errors.NewWithCause("Shutting down the server returning error", err))
		}
	}
}

// CommonHTTPMiddlware .
func CommonHTTPMiddlware(withOutLog bool) []func(http.HandlerFunc) http.HandlerFunc {
	loggerMiddleware := middleware.Nop
	if withOutLog {
		loggerMiddleware = logger.NewWithDefaultLogger(Out)
	}
	return []func(http.HandlerFunc) http.HandlerFunc{
		paniclogger.New(0, func(ev paniclogger.Event) {
			if err, ok := ev.Error.(error); ok {
				EprintTime(errors.Wrap(err))
			} else {
				EprintTime(errors.Errorf("not an error panic: %v", ev.Error))
			}
		}),
		kv.New(),
		loggerMiddleware,
	}
}
