package log4g

import "time"

type Level int

const levelStep = 10
const (
	FATAL Level = levelStep*iota + levelStep
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
	ALL
)

type Logger interface {
	Fatal(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
	Logf(level Level, fstr string, args ...interface{})
	Log(level Level, payload interface{})
}

type LogEvent struct {
	Level      Level
	Timestamp  time.Time
	LoggerName string
	Payload    interface{}
}

type Appender interface {
	Append(event *LogEvent) bool
}

// The factory allows to create an appender instances
type AppenderFactory interface {
	// Appender name
	Name() string
	NewAppender(map[string]interface{}) (Appender, error)
	Shutdown()
}

/**
 * Provides pointer to the logger with specified name.
 * name can have 'dot' separated form.
 */
func GetLogger(name string) Logger {
	return lm.getLogger(name)
}

/**
 * Returns slice with log level names. Changing the appropriate level name here will
 * follow to changing its name in log messages for appenders that form the message
 * from LogEvent
 */
func LevelNames() []string {
	return lm.levelNames
}

/**
 * All appenders should register them in their module init() method.
 * The method returns error if the function is called after config intialization sub-system.
 * Parameters:
 *		appenderFactory - interface which allows to create new instances of
 * 			some specific appender type.
 */
func RegisterAppender(appenderFactory AppenderFactory) error {
	return lm.registerAppender(appenderFactory)
}

/**
 * Should be called to shutdown log subsystem properly. It will notify all logContexts and wait
 * while all go routines are over. To call this method could be essential to finalize some
 * appenders implementations and close them properly
 */
func Shutdown() {
	// TODO: implement it
}
