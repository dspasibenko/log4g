package log4g

import (
	"errors"
	"github.com/dspasibenko/log4g/collections"
	"sync"
)

const rootLoggerName = ""

type logManager struct {
	loggers     map[string]*logger
	logLevels   *collections.SortedSlice
	logContexts *collections.SortedSlice
	appenders   map[string]*appenderHolder
	rwLock      sync.RWMutex
	levelNames  []string
}

type appenderHolder struct {
	newFunc NewAppenderFn
}

var lm logManager

func init() {
	lm.logLevels, _ = collections.NewSortedSlice(5)
	rootLLS := &logLevelSetting{rootLoggerName, INFO}
	lm.logLevels.Add(rootLLS)

	lm.loggers = make(map[string]*logger)
	lm.loggers[rootLoggerName] = &logger{rootLoggerName, rootLLS, rootLLS.level}
	lm.logContexts, _ = collections.NewSortedSlice(2)
	lm.appenders = make(map[string]*appenderHolder)

	lm.levelNames = make([]string, ALL+1)
	lm.levelNames[FATAL] = "FATAL"
	lm.levelNames[ERROR] = "ERROR"
	lm.levelNames[WARN] = "WARN "
	lm.levelNames[INFO] = "INFO "
	lm.levelNames[DEBUG] = "DEBUG"
	lm.levelNames[TRACE] = "TRACE"
	lm.levelNames[ALL] = "ALL  "
}

func (lm *logManager) getLogger(loggerName string) Logger {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	loggerName = normalizeLogName(loggerName)
	l, ok := lm.loggers[loggerName]
	if !ok {
		// Create new logger for the name
		rootLLS := getLogLevelSetting(loggerName, lm.logLevels)
		l = &logger{loggerName, rootLLS, rootLLS.level}
		lm.loggers[loggerName] = l
	}
	return l
}

func (lm *logManager) setLogLevel(loggerName string, level Level) {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	lls := setLogLevel(level, loggerName, lm.logLevels)
	applyNewLevelToLoggers(lls, lm.loggers)
}

func (lm *logManager) registerAppender(appenderName string, newAppenderFn NewAppenderFn) error {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	appenderName = normalizeLogName(appenderName)
	_, ok := lm.appenders[appenderName]
	if ok {
		return errors.New("Cannot register appender with the name " + appenderName +
			" because appender with the name is already registerd ")
	}

	lm.appenders[appenderName] = &appenderHolder{newFunc: newAppenderFn}
	return nil
}
