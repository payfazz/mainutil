package mainutil

import (
	"github.com/payfazz/stdlog"
)

// Env .
type Env struct {
	Info *stdlog.Logger
	Err  *stdlog.Logger
}

// InfoLogger .
func (env *Env) InfoLogger() *stdlog.Logger {
	if env != nil && env.Info != nil {
		return env.Info
	}

	return stdlog.Out
}

// ErrLogger .
func (env *Env) ErrLogger() *stdlog.Logger {
	if env != nil && env.Err != nil {
		return env.Err
	}

	return stdlog.Err
}
