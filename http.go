package mainutil

import (
	"context"
	"log"
	"net"
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
	"github.com/payfazz/stdlog"
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

	waitHTTP(serverErrCh, server)
}

// RunHTTPServerOn .
func RunHTTPServerOn(server *http.Server, listener net.Listener) {
	serverErrCh := make(chan error, 1)
	go func() {
		defer close(serverErrCh)
		if server.TLSConfig == nil {
			Iprintf("Server listen on \"%s\"\n", listener.Addr())
			serverErrCh <- errors.Wrap(server.Serve(listener))
		} else {
			Iprintf("Server listen on TLS \"%s\"\n", listener.Addr())
			serverErrCh <- errors.Wrap(server.ServeTLS(listener, "", ""))
		}
	}()

	waitHTTP(serverErrCh, server)
}

func waitHTTP(serverErrCh chan error, server *http.Server) {
	signalChan := make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	signal.Notify(signalChan, signals...)

	select {
	case err := <-serverErrCh:
		signal.Reset(signals...)
		errhandler.Check(errors.NewWithCause("listener goroutine got an error", errors.Wrap(err)))
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
			errhandler.Check(errors.NewWithCause("Shutting down the server returning error", errors.Wrap(err)))
		}
	}
}

// CommonHTTPMiddlware .
func CommonHTTPMiddlware(haveOutLog bool) []func(http.HandlerFunc) http.HandlerFunc {
	loggerMiddleware := middleware.Nop
	if haveOutLog {
		loggerMiddleware = logger.NewWithDefaultLogger(Out)
	}
	return []func(http.HandlerFunc) http.HandlerFunc{
		paniclogger.New(0, func(ev paniclogger.Event) {
			if err, ok := ev.Error.(error); ok {
				errors.PrintTo(Err, errors.Wrap(err))
			} else {
				errors.PrintTo(Err, errors.Errorf("not an error panic: %v", ev.Error))
			}
		}),
		kv.New(),
		loggerMiddleware,
	}
}

// DefaultHTTPServer .
func DefaultHTTPServer() *http.Server {
	ret := http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		ErrorLog:          log.New(stdlog.Err, "internal http error: ", log.LstdFlags|log.LUTC),
	}
	return &ret
}
