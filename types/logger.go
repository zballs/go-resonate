package types

import (
	log "github.com/inconshreveable/log15"
	. "github.com/zballs/go_resonate/util"
)

type Logger struct {
	log.Logger
}

func NewLogger(module string) Logger {
	return Logger{log.New("module", module)}
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.Error(Sprintf(format, args...))
}
