package log4g

import (
	"fmt"
	"time"
)

type logger struct {
	loggerName string
	lls        *logLevelSetting
	lctx       *logContext
	logLevel   Level
}

func (l *logger) Fatal(args ...interface{}) {
	fmt.Println(args)
}

func (l *logger) Error(args ...interface{}) {

}

func (l *logger) Warn(args ...interface{}) {

}

func (l *logger) Info(args ...interface{}) {

}

func (l *logger) Debug(args ...interface{}) {

}

func (l *logger) Trace(args ...interface{}) {

}

func (l *logger) Logf(level Level, fstr string, args ...interface{}) {

}

func (l *logger) Log(level Level, payload interface{}) {
	if l.logLevel < level {
		return
	}
	l.lctx.log(&LogEvent{level, time.Now(), l.loggerName, payload})
}

func (l *logger) setLogLevelSetting(lls *logLevelSetting) {
	l.lls = lls
	l.logLevel = lls.level
}

func (l *logger) setLogContext(lctx *logContext) {
	l.lctx = lctx
}

// Apply new LogLevelSetting to all appropriate loggers
func applyNewLevelToLoggers(lls *logLevelSetting, loggers map[string]*logger) {
	for _, l := range loggers {
		if !ancestor(lls.loggerName, l.loggerName) {
			continue
		}
		if ancestor(l.lls.loggerName, lls.loggerName) {
			l.setLogLevelSetting(lls)
		}
	}
}
