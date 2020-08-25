package mainutil

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/payfazz/go-errors"
	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/logger"
	"github.com/payfazz/go-middleware/common/paniclogger"
	"github.com/payfazz/stdlog"
)

func httpSetTLSInternal(s *http.Server, tls *tls.Config) {
	tls.NextProtos = []string{"h2", "http/1.1"}
	s.TLSConfig = tls
}

// HTTPSetTLS .
func HTTPSetTLS(s *http.Server, certfile string, keyfile string) error {
	tls, err := DefaultTLSConfig(certfile, keyfile)
	if err != nil {
		return errors.Wrap(err)
	}
	httpSetTLSInternal(s, tls)
	return nil
}

// HTTPSetTLSString .
func HTTPSetTLSString(s *http.Server, certpem string, keypem string) error {
	tls, err := DefaultTLSConfigString(certpem, keypem)
	if err != nil {
		return errors.Wrap(err)
	}
	httpSetTLSInternal(s, tls)
	return nil
}

// DefaultHTTPServer .
func DefaultHTTPServer(addr string, handler http.HandlerFunc) *http.Server {
	s := &http.Server{}

	s.ReadTimeout = 1 * time.Minute
	s.WriteTimeout = 1 * time.Minute
	s.IdleTimeout = 30 * time.Second

	s.ErrorLog = stdlog.NewFromEnv(stdlog.Err(), "net/http.Server.ErrorLog: ").AsLogger()

	s.Addr = addr
	s.Handler = handler

	return s
}

// CommonHTTPMiddlware .
func CommonHTTPMiddlware(printRequestLog bool) []func(http.HandlerFunc) http.HandlerFunc {
	reqLoggerMiddleware := middleware.Nop
	if printRequestLog {
		reqLoggerMiddleware = logger.NewWithDefaultLogger(stdlog.Out())
	}

	return []func(http.HandlerFunc) http.HandlerFunc{
		paniclogger.New(0, func(ev paniclogger.Event) { printErr(ev.Error) }),
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
		serverAddr := s.Addr
		if l != nil {
			serverAddr = l.Addr().String()
		}
		stdlog.PrintOut(fmt.Sprintf(
			"Shutting down the server \"%s\" (Waiting for graceful shutdown: %s)\n",
			serverAddr,
			gracefulShutdown.Truncate(time.Second).String(),
		))
		if err := s.Shutdown(shutdownCtx); err != nil {
			errors.PrintTo(stdlog.Err(),
				errors.NewWithCause("fail to graceful shutdown", err),
			)
		}
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
