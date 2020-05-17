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
func ListenUnixSocket(path string) (listener net.Listener, cleanUpFunc func(context.Context), err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}

	listener, err = net.Listen("unix", absPath)
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}

	cleanUpFunc = func(ctx context.Context) {
		<-ctx.Done()
		os.Remove(absPath)
		listener.Close()
	}

	return listener, cleanUpFunc, nil
}
