package log4g

import (
	"errors"
	"sync"
)

type logManager struct {
	config *logConfig
	rwLock sync.RWMutex
}

var lm logManager

func init() {
	lm.config = newLogConfig()
}

func (lm *logManager) getLogger(loggerName string) Logger {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	lm.config.initIfNeeded()
	return lm.config.getLogger(loggerName)
}

func (lm *logManager) setLogLevel(loggerName string, level Level) {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	lm.config.initIfNeeded()
	lm.config.setLogLevel(level, loggerName)
}

func (lm *logManager) registerAppender(appenderFactory AppenderFactory) error {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	appenderName := appenderFactory.Name()
	_, ok := lm.config.appenderFactorys[appenderName]
	if ok {
		return errors.New("Cannot register appender factory for the name " + appenderName +
			" because the name is already registerd ")
	}

	lm.config.appenderFactorys[appenderName] = appenderFactory
	return nil
}
