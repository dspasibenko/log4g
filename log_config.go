package log4g

import (
	"github.com/dspasibenko/log4g/collections"
	"strconv"
	"strings"
)

type logConfig struct {
	loggers          map[string]*logger
	logLevels        *collections.SortedSlice
	logContexts      *collections.SortedSlice
	appenderFactorys map[string]AppenderFactory
	appenders        map[string]Appender
	levelNames       []string
	levelMap         map[string]Level
}

// Config params
const (
	// appender.console.type=log4g/appenders/consoleAppender
	cfgAppender     = "appender"
	cfgAppenderType = "type"

	//level.11=SEVERE
	cfgLevel = "level"

	// context.a.b.c.appenders=console,ROOT
	// context.a.b.c.level=INFO
	// context.a.b.c.buffer=100
	cfgContext          = "context"
	cfgContextAppenders = "appenders"
	cfgContextLevel     = "level"
	cfgContextBufSize   = "buffer"
	cfgContextInherited = "inherited"

	// logger.a.b.c.d.level=INFO
	cfgLogger      = "logger"
	cfgLoggerLevel = "level"
)

const rootLoggerName = ""

var defaultConfigParams = map[string]string{
	"appender.ROOT.layout": "[%d{01-02 15:04:05.000}] %p %c: %m",
	"appender.ROOT.type":   "log4g/appenders/consoleAppender",
	"context.appenders":    "ROOT",
	"context.level":        "INFO"}

func newLogConfig() *logConfig {
	lc := &logConfig{}

	lc.loggers = make(map[string]*logger)
	lc.logLevels, _ = collections.NewSortedSlice(10)
	lc.logContexts, _ = collections.NewSortedSlice(2)
	lc.appenderFactorys = make(map[string]AppenderFactory)
	lc.appenders = make(map[string]Appender)
	lc.levelNames = make([]string, ALL+1)
	lc.levelMap = make(map[string]Level)

	lc.levelNames[FATAL] = "FATAL"
	lc.levelNames[ERROR] = "ERROR"
	lc.levelNames[WARN] = "WARN "
	lc.levelNames[INFO] = "INFO "
	lc.levelNames[DEBUG] = "DEBUG"
	lc.levelNames[TRACE] = "TRACE"
	lc.levelNames[ALL] = "ALL  "

	return lc
}

func mergedParamsWithDefault(params map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range params {
		result[k] = v
	}
	for k, v := range defaultConfigParams {
		_, ok := result[k]
		if !ok {
			result[k] = v
		}
	}
	return result
}

// Check whether ROOT context exists, if no, it will be initialized by default
func (lc *logConfig) initIfNeeded() {
	if getLogLevelContext(rootLoggerName, lc.logContexts) != nil {
		return
	}
	lc.setConfigParams(defaultConfigParams)
}

func (lc *logConfig) initWithParams(oldLogConfig *logConfig, params map[string]string) {
	for k, v := range oldLogConfig.loggers {
		lc.loggers[k] = v
	}

	for k, v := range oldLogConfig.appenderFactorys {
		lc.appenderFactorys[k] = v
	}

	lc.logLevels, _ = collections.NewSortedSliceByParams(oldLogConfig.logLevels.Copy()...)
	lc.setConfigParams(params)
}

func (lc *logConfig) setConfigParams(params map[string]string) {
	lc.applyLevelParams(params)
	lc.createAppenders(params)
	lc.createContexts(params)
	lc.createLoggers(params)
	lc.applyLevels()
}

// Allows to specify custom level names in form level.X=<levelName>
// for example: level.11=SEVERE
func (lc *logConfig) applyLevelParams(params map[string]string) {
	for i, _ := range lc.levelNames {
		param := cfgLevel + strconv.Itoa(i)
		v, ok := params[param]
		if ok {
			lc.levelNames[i] = v
		}
		level := strings.Trim(strings.ToLower(lc.levelNames[i]), " ")
		if len(level) > 0 {
			lc.levelMap[level] = Level(i)
		}
	}
}

func (lc *logConfig) createAppenders(params map[string]string) {
	// collect settings for all appenders from config
	apps := groupConfigParams(params, cfgAppender)

	// create appenders
	for appName, appAttributes := range apps {
		t := appAttributes[cfgAppenderType]
		f, ok := lc.appenderFactorys[t]
		if !ok {
			panic("No Factory for appender " + t)
		}

		app, err := f.NewAppender(appAttributes)
		if err != nil {
			panic(err)
		}

		lc.appenders[appName] = app
	}
}

func (lc *logConfig) createContexts(params map[string]string) {
	// collect settings for all contexts from config
	ctxs := groupConfigParams(params, cfgContext)

	// create contexts
	for logName, ctxAttributes := range ctxs {
		appenders := lc.getAppendersFromList(ctxAttributes[cfgContextAppenders])
		if len(appenders) == 0 {
			panic("context " + logName + " doesn't refer at least to one appender")
		}

		level := lc.getLevelByName(ctxAttributes[cfgContextLevel])
		bufSizeStr, ok := ctxAttributes[cfgContextBufSize]
		bufSize := 100
		if ok {
			bufSize, err := strconv.Atoi(strings.Trim(bufSizeStr, " "))
			if err != nil || bufSize < 0 {
				panic("Incorrect buffer size value for context attribute " + cfgContextBufSize + " should be positive integer")
			}
		}

		inhStr, ok := ctxAttributes[cfgContextInherited]
		inh := true
		if ok {
			var err error
			inh, err = strconv.ParseBool(inhStr)
			if err != nil {
				panic("Incorrect context attibute " + cfgContextInherited + " value, should be true or false")
			}
		}

		setLogLevel(level, logName, lc.logLevels)
		context, _ := newLogContext(logName, appenders, inh, bufSize)
		lc.logContexts.Add(context)
	}
}

func (lc *logConfig) createLoggers(params map[string]string) {
	// collect settings for all loggers from config
	loggers := groupConfigParams(params, cfgLogger)

	// apply logger settings
	for loggerName, loggerAttributes := range loggers {
		level := lc.getLevelByName(loggerAttributes[cfgLoggerLevel])
		setLogLevel(level, loggerName, lc.logLevels)
		lc.getLogger(loggerName)
	}
}

func (lc *logConfig) applyLevels() {
	for _, l := range lc.loggers {
		rootLLS := getLogLevelSetting(l.loggerName, lc.logLevels)
		l.setLogLevelSetting(rootLLS)

		lctx := getLogLevelContext(l.loggerName, lc.logContexts)
		l.setLogContext(lctx)
	}
}

func (lc *logConfig) getLogger(loggerName string) Logger {
	loggerName = normalizeLogName(loggerName)
	l, ok := lc.loggers[loggerName]
	if !ok {
		// Create new logger for the name
		rootLLS := getLogLevelSetting(loggerName, lc.logLevels)
		rootCtx := getLogLevelContext(loggerName, lc.logContexts)
		l = &logger{loggerName, rootLLS, rootCtx, rootLLS.level}
		lc.loggers[loggerName] = l
	}
	return l
}

func (lc *logConfig) setLogLevel(level Level, loggerName string) {
	lls := setLogLevel(level, loggerName, lc.logLevels)
	applyNewLevelToLoggers(lls, lc.loggers)
}

func (lc *logConfig) getAppendersFromList(appNames string) []Appender {
	names := strings.Split(appNames, ",")
	result := make([]Appender, 0, len(names))
	for _, name := range names {
		name = strings.Trim(name, " ")
		a, ok := lc.appenders[name]
		if !ok {
			continue
		}
		result = append(result, a)
	}
	return result
}

// gets level index, or -1 if not found
func (lc *logConfig) getLevelByName(levelName string) (idx Level) {
	levelName = strings.Trim(strings.ToLower(levelName), " ")
	idx, ok := lc.levelMap[levelName]
	if !ok {
		idx = -1
	}
	return idx
}
