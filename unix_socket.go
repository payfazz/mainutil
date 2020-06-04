// +build linux

package mainutil

import (
	"net"
	"os"
	"path/filepath"

	"github.com/payfazz/go-errors"
)

// ListenUnixSocket .
func ListenUnixSocket(path string) (listener net.Listener, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	listener, err = net.Listen("unix", absPath)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &unixSocketWrapper{absPath: absPath, Listener: listener}, nil
}

type unixSocketWrapper struct {
	absPath string
	net.Listener
}

func (u *unixSocketWrapper) Close() error {
	os.RemoveAll(u.absPath)
	return u.Listener.Close()
}
