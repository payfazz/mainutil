// +build linux

package mainutil

import (
	"context"
	"net"
	"os"

	"github.com/payfazz/go-errors"
)

// ListenUnixSocket .
func (env *Env) ListenUnixSocket(ctx context.Context, path string) (net.Listener, error) {
	l, err := net.Listen("unix", path)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	go func() {
		<-ctx.Done()
		os.RemoveAll(path)
	}()

	return l, nil
}
