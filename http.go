package mainutil

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/payfazz/go-errors"
)

// RunHTTPServer run *http.Server,
// when SIGTERM or SIGINT is recieved graceful shutdown the server.
//
// github.com/payfazz/go-errors.HandleWith must be already defered.
func RunHTTPServer(server *http.Server) {
	serverErrCh := make(chan error, 1)
	go func() {
		defer close(serverErrCh)
		if server.TLSConfig == nil {
			Iprintf("Server listen on \"%s\"\n", server.Addr)
			serverErrCh <- server.ListenAndServe()
		} else {
			Iprintf("Server listen on TLS \"%s\"\n", server.Addr)
			serverErrCh <- server.ListenAndServeTLS("", "")
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	signal.Notify(signalChan, signals...)

	select {
	case err := <-serverErrCh:
		errors.Fail("Server returning error", err)
	case sig := <-signalChan:
		signal.Reset(signals...)
		waitFor := 1 * time.Minute
		Iprintf(
			"Got '%s' signal, Stopping (Waiting for graceful shutdown: %s)\n",
			sig.String(), waitFor.String(),
		)
		ctx, cancel := context.WithTimeout(context.Background(), waitFor)
		defer cancel()
		errors.CheckOrFail("Shutting down server returning error", server.Shutdown(ctx))
	}
}
