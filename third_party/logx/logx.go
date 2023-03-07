package logx

import (
	"github.com/go-kratos/kratos/v2/log"
)

type Option func(logger *filterLogger)
type filterLogger struct {
	logger log.Logger
	level  log.Level
}

func WithFilterLevel(level log.Level) func(logger *filterLogger) {
	return func(logger *filterLogger) {
		logger.level = level
	}
}

// NewLogger new a logger with logger and filter level.
func NewLogger(logger log.Logger, opts ...Option) log.Logger {
	l := &filterLogger{
		logger: logger,
		level:  log.LevelDebug,
	}
	for _, o := range opts {
		o(l)
	}
	return l

}

func (fl *filterLogger) Log(level log.Level, keyvals ...interface{}) error {
	if level < fl.level {
		return nil
	}
	return fl.logger.Log(level, keyvals...)
}
