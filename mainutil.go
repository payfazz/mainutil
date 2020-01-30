package mainutil

import (
	"github.com/payfazz/stdlog"
)

// Env .
type Env struct {
	Info *stdlog.Logger
	Err  *stdlog.Logger
}

func (env *Env) info() *stdlog.Logger {
	if env != nil && env.Info != nil {
		return env.Info
	}

	return stdlog.Out
}

func (env *Env) err() *stdlog.Logger {
	if env != nil && env.Err != nil {
		return env.Err
	}

	return stdlog.Err
}
