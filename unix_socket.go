// +build linux

package mainutil

import (
	"context"
	"net"
	"os"
	"path/filepath"

	"github.com/payfazz/go-errors"
)

// ListenUnixSocket .
func (env *Env) ListenUnixSocket(path string) (net.Listener, func(context.Context), error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}

	l, err := net.Listen("unix", absPath)
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}

	cleanUpFunc := func(ctx context.Context) {
		<-ctx.Done()
		l.Close()
		os.RemoveAll(absPath)
	}

	return l, cleanUpFunc, nil
}
