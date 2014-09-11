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
	Log(level Level, args ...interface{})
	Logf(level Level, fstr string, args ...interface{})
	Logp(level Level, payload interface{})
}

type LogEvent struct {
	Level      Level
	Timestamp  time.Time
	LoggerName string
	Payload    interface{}
}

type Appender interface {
	Append(event *LogEvent) bool
	// should be called every time when the instance is not going to be used anymore
	Shutdown()
}

// The factory allows to create an appender instances
type AppenderFactory interface {
	// Appender name
	Name() string
	NewAppender(map[string]string) (Appender, error)
	Shutdown()
}

/**
 * Provides pointer to the logger with specified name.
 * name can have 'dot' separated form.
 */
func GetLogger(loggerName string) Logger {
	return lm.getLogger(loggerName)
}

func SetLogLevel(loggerName string, level Level) {
	lm.setLogLevel(loggerName, level)
}

/**
 * Returns slice with log level names. Changing the appropriate level name here will
 * follow to changing its name in log messages for appenders that form the message
 * by provided LogEvent values.
 */
func LevelNames() []string {
	return lm.config.levelNames
}

/**
 * All appenders should register them in their module init() method or by calling this function directly.
 * The method returns error if the function is called after config intialization sub-system.
 * Parameters:
 *		appenderFactory - interface which allows to create new instances of
 * 			some specific appender type.
 */
func RegisterAppender(appenderFactory AppenderFactory) error {
	return lm.registerAppender(appenderFactory)
}

/**
 * Reads log4g configuration properties from text file, which name is provided as
 * configFileName parameter.
 */
func ConfigF(configFileName string) error {
	return lm.setPropsFromFile(configFileName)
}

/**
 * Configures log4g by key:value pairs provided as a map of properties
 */
func Config(props map[string]string) error {
	return lm.setNewProperties(props)
}

/**
 * Should be called to shutdown log subsystem properly. It will notify all logContexts and wait
 * while all go routines are over. To call this method could be essential to finalize some
 * appenders implementations and close them properly
 */
func Shutdown() {
	lm.shutdown()
}
