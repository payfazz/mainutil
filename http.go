package mainutil

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/payfazz/go-errors"
	"github.com/payfazz/go-errors/errhandler"
	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/logger"
	"github.com/payfazz/go-middleware/common/paniclogger"
)

// SetDefaultForHTTP .
func (env *Env) SetDefaultForHTTP(s *http.Server) {
	s.ReadTimeout = 1 * time.Minute
	s.WriteTimeout = 1 * time.Minute
	s.IdleTimeout = 30 * time.Second
	s.ErrorLog = log.New(env.ErrLogger(), "internal http error: ", log.LstdFlags|log.LUTC)
}

// DefaultHTTPServer .
func (env *Env) DefaultHTTPServer(addr string, handler http.HandlerFunc) *http.Server {
	s := http.Server{}
	env.SetDefaultForHTTP(&s)
	s.Addr = addr
	s.Handler = handler
	return &s
}

// CommonHTTPMiddlware .
func (env *Env) CommonHTTPMiddlware(haveOutLog bool) []func(http.HandlerFunc) http.HandlerFunc {
	requestLogger := middleware.Nop
	if haveOutLog {
		requestLogger = logger.NewWithDefaultLogger(env.InfoLogger())
	}

	errLogger := log.New(env.ErrLogger(), "unhandled panic: ", log.LstdFlags|log.LUTC)

	return []func(http.HandlerFunc) http.HandlerFunc{
		paniclogger.New(0, func(ev paniclogger.Event) {
			if err, ok := ev.Error.(error); ok {
				errors.PrintTo(errLogger, errors.Wrap(errhandler.UnwrapUnhandledError(err)))
			} else {
				errors.PrintTo(errLogger, errors.Errorf("non error panic: %v", ev.Error))
			}
		}),
		requestLogger,
	}
}

// RunHTTPServerOn .
func (env *Env) RunHTTPServerOn(
	ctx context.Context,
	s *http.Server,
	l net.Listener,
	gracefulShutdown time.Duration,
) error {
	serverErrCh := make(chan error, 1)
	go func() {
		defer close(serverErrCh)
		if l == nil {
			serverErrCh <- env.runHTTPServerOnDefaultListener(s)
		} else {
			serverErrCh <- env.runHTTPServerOnListener(s, l)
		}
	}()

	select {
	case err := <-serverErrCh:
		return errors.Wrap(err)
	case <-ctx.Done():
		if gracefulShutdown == 0 {
			gracefulShutdown = max(s.ReadTimeout, s.WriteTimeout)
		}
		if gracefulShutdown == 0 {
			gracefulShutdown = 1 * time.Minute
		}
		gracefulShutdown += 500 * time.Millisecond
		shutdownCtx, cancel := context.WithTimeout(ctx, gracefulShutdown)
		defer cancel()
		env.InfoLogger().Print(fmt.Sprintf(
			"Shutting down the server (Waiting for graceful shutdown: %s)\n",
			gracefulShutdown.Truncate(time.Second).String(),
		))
		return errors.Wrap(s.Shutdown(shutdownCtx))
	}
}

func (env *Env) runHTTPServerOnDefaultListener(s *http.Server) error {
	if s.TLSConfig != nil {
		env.InfoLogger().Print(fmt.Sprintf("Server listen on TLS \"%s\"\n", s.Addr))
		return errors.Wrap(s.ListenAndServeTLS("", ""))
	}

	env.InfoLogger().Print(fmt.Sprintf("Server listen on \"%s\"\n", s.Addr))
	return errors.Wrap(s.ListenAndServe())
}

func (env *Env) runHTTPServerOnListener(s *http.Server, l net.Listener) error {
	if s.TLSConfig != nil {
		env.InfoLogger().Print(fmt.Sprintf("Server listen on TLS \"%s\"\n", l.Addr().String()))
		return errors.Wrap(s.ServeTLS(l, "", ""))
	}

	env.InfoLogger().Print(fmt.Sprintf("Server listen on \"%s\"\n", l.Addr().String()))
	return errors.Wrap(s.Serve(l))
}

func max(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
