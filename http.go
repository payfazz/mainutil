package mainutil

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/payfazz/go-errors"
	"github.com/payfazz/go-errors/errhandler"
	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
	"github.com/payfazz/go-middleware/common/logger"
	"github.com/payfazz/go-middleware/common/paniclogger"
	"github.com/payfazz/stdlog"
)

// HTTPSetDefault .
func HTTPSetDefault(s *http.Server) {
	s.ReadTimeout = 1 * time.Minute
	s.WriteTimeout = 1 * time.Minute
	s.IdleTimeout = 30 * time.Second
	s.ErrorLog = stdlog.NewFromEnv(stdlog.Err(), "net/http.Server.ErrorLog: ").AsLogger()
}

// HTTPSetTLS .
func HTTPSetTLS(s *http.Server, certfile string, keyfile string) error {
	tls, err := DefaultTLSConfig(certfile, keyfile)
	if err != nil {
		return errors.Wrap(err)
	}
	tls.NextProtos = []string{"h2", "http/1.1"}

	s.TLSConfig = tls

	return nil
}

// DefaultHTTPServer .
func DefaultHTTPServer(addr string, handler http.HandlerFunc) *http.Server {
	s := http.Server{}
	HTTPSetDefault(&s)
	s.Addr = addr
	s.Handler = handler
	return &s
}

// CommonHTTPMiddlware .
func CommonHTTPMiddlware(printRequestLog bool) []func(http.HandlerFunc) http.HandlerFunc {
	reqLoggerMiddleware := middleware.Nop
	if printRequestLog {
		reqLoggerMiddleware = logger.NewWithDefaultLogger(stdlog.Out())
	}

	return []func(http.HandlerFunc) http.HandlerFunc{
		paniclogger.New(0, func(ev paniclogger.Event) {
			if err, ok := ev.Error.(error); ok {
				errors.PrintTo(stdlog.Err(), errors.Wrap(errhandler.UnwrapUnhandledError(err)))
			} else {
				errors.PrintTo(stdlog.Err(), errors.Errorf("non-error-type: %v", ev.Error))
			}
		}),
		kv.New(),
		reqLoggerMiddleware,
	}
}

// RunHTTPServerOn .
func RunHTTPServerOn(
	ctx context.Context,
	s *http.Server,
	l net.Listener,
	gracefulShutdown time.Duration,
) error {
	serverErrCh := make(chan error, 1)
	go func() {
		defer close(serverErrCh)
		if l == nil {
			serverErrCh <- errors.Wrap(runHTTPServerOnDefaultListener(s))
		} else {
			serverErrCh <- errors.Wrap(runHTTPServerOnListener(s, l))
		}
	}()

	select {
	case err := <-serverErrCh:
		if err != nil {
			return errors.NewWithCause("Fail to run server", err)
		}
		return nil
	case <-ctx.Done():
		maxDuration := func(a, b time.Duration) time.Duration {
			if a > b {
				return a
			}
			return b
		}

		if gracefulShutdown == 0 {
			gracefulShutdown = maxDuration(s.ReadTimeout, s.WriteTimeout)
		}
		if gracefulShutdown == 0 {
			gracefulShutdown = 1 * time.Minute
		}
		gracefulShutdown += 500 * time.Millisecond // give more 0.5 second for cleanup
		shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdown)
		defer cancel()
		stdlog.PrintOut(fmt.Sprintf(
			"Shutting down the server (Waiting for graceful shutdown: %s)\n",
			gracefulShutdown.Truncate(time.Second).String(),
		))
		s.Shutdown(shutdownCtx)
		return nil
	}
}

func runHTTPServerOnDefaultListener(s *http.Server) error {
	if s.TLSConfig != nil {
		stdlog.PrintOut(fmt.Sprintf("Server listen on TLS \"%s\"\n", s.Addr))
		if err := s.ListenAndServeTLS("", ""); err != nil {
			return errors.Wrap(err)
		}
		return nil
	}

	stdlog.PrintOut(fmt.Sprintf("Server listen on \"%s\"\n", s.Addr))
	if err := s.ListenAndServe(); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func runHTTPServerOnListener(s *http.Server, l net.Listener) error {
	if s.TLSConfig != nil {
		stdlog.PrintOut(fmt.Sprintf("Server listen on TLS \"%s\"\n", l.Addr().String()))
		if err := s.ServeTLS(l, "", ""); err != nil {
			errors.Wrap(err)
		}
		return nil
	}

	stdlog.PrintOut(fmt.Sprintf("Server listen on \"%s\"\n", l.Addr().String()))
	if err := s.Serve(l); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
