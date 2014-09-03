package log4g

import "sync"

type logManager struct {
	loggers   map[string]*logger
	logLevels map[string]Level
	rwLock    sync.RWMutex
}

var lm logManager

func init() {
	lm.loggers = make(map[string]*logger)
}

func (lm *logManager) getLogger(name string) Logger {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	l, ok := lm.loggers[name]
	if !ok {
		// Create new logger for the name
		l = new(logger)
		lm.loggers[name] = l
	}
	return l
}
