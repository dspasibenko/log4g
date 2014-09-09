package log4g

import (
	"errors"
	"github.com/dspasibenko/log4g/collections"
	"strconv"
	"strings"
	"sync"
)

const rootLoggerName = ""

type logManager struct {
	loggers          map[string]*logger
	logLevels        *collections.SortedSlice
	logContexts      *collections.SortedSlice
	appenderFactorys map[string]AppenderFactory
	appenders        map[string]Appender
	rwLock           sync.RWMutex
	levelNames       []string
	levelMap         map[string]Level
}

var lm logManager

func init() {
	lm.logLevels, _ = collections.NewSortedSlice(5)
	rootLLS := &logLevelSetting{rootLoggerName, INFO}
	lm.logLevels.Add(rootLLS)

	lm.loggers = make(map[string]*logger)
	lm.loggers[rootLoggerName] = &logger{rootLoggerName, rootLLS, rootLLS.level}
	lm.logContexts, _ = collections.NewSortedSlice(2)
	lm.appenderFactorys = make(map[string]AppenderFactory)
	lm.appenders = make(map[string]Appender)

	lm.levelNames = make([]string, ALL+1)
	lm.levelNames[FATAL] = "FATAL"
	lm.levelNames[ERROR] = "ERROR"
	lm.levelNames[WARN] = "WARN "
	lm.levelNames[INFO] = "INFO "
	lm.levelNames[DEBUG] = "DEBUG"
	lm.levelNames[TRACE] = "TRACE"
	lm.levelNames[ALL] = "ALL  "

	// Apply default params
	// TODO: not here!
	/*	lm.setConfig(map[string]string{
		"appender.ROOT.layout": "[%d{01-02 15:04:05.000}] %p %c: %m",
		"appender.ROOT.type":   "log4g/appenders/consoleAppender",
		"context.appenders":    "ROOT",
		"context.level":        "INFO",
		"context.inherited":    "true",
		"level.10":             "FATAL",
		"level.20":             "ERROR",
		"level.30":             "WARN ",
		"level.40":             "INFO ",
		"level.50":             "DEBUG",
		"level.60":             "TRACE",
		"level.70":             "ALL  ",
	}) */
}

// --------------- Config functions ---------------
// expects list of properties and apply it to current configuration
func (lm *logManager) setConfig(params map[string]string) {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	lm.applyLevelParams(params)
	lm.createAppenders(params)
	// parse context
	// parse loggers
}

func (lm *logManager) applyLevelParams(params map[string]string) {
	lm.levelMap = make(map[string]Level)
	for i, _ := range lm.levelNames {
		param := cfgLevel + strconv.Itoa(i)
		v, ok := params[param]
		if ok {
			lm.levelNames[i] = v
		}
		level := strings.Trim(strings.ToLower(lm.levelNames[i]), " ")
		if len(level) > 0 {
			lm.levelMap[level] = Level(i)
		}
	}
}

func (lm *logManager) createAppenders(params map[string]string) {
	// collect settings for all appenders from config
	apps := parseAppendersParams(params)

	// create appenders
	for appName, appParams := range apps {
		t := appParams[cfgAppenderType]
		f, ok := lm.appenderFactorys[t]
		if !ok {
			panic("No Factory for appender " + t)
		}

		app, err := f.NewAppender(appParams)
		if err != nil {
			panic(err)
		}

		lm.appenders[appName] = app
	}
}

func (lm *logManager) createContexts(params map[string]string) {
	// collect settings for all contexts from config
	//ctxts := parseContextParams(params)
	//for logName, ctxParams := range params {

	//}
}

// --------------- Other getters ----------------
func (lm *logManager) getLevelByName(levelName string) (idx Level, ok bool) {
	levelName = strings.Trim(strings.ToLower(levelName), " ")
	idx, ok = lm.levelMap[levelName]
	return idx, ok
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

func (lm *logManager) registerAppender(appenderFactory AppenderFactory) error {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	appenderName := normalizeLogName(appenderFactory.Name())
	_, ok := lm.appenderFactorys[appenderName]
	if ok {
		return errors.New("Cannot register appender factory for the name " + appenderName +
			" because the name is already registerd ")
	}

	lm.appenderFactorys[appenderName] = appenderFactory
	return nil
}
