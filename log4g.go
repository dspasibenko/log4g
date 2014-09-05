package log4g

type Level int

const (
	FATAL Level = iota
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
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

type Timestamp uint64

type LogEvent struct {
	level     Level
	timestamp Timestamp
	logger    string
	payload   interface{}
}

type Appender interface {
	Append(record *LogEvent)
}

type NewAppenderFn func(map[string]interface{}) Appender

/**
 * Provides pointer to the logger with specified name.
 * name can have 'dot' separated form.
 */
func GetLogger(name string) Logger {
	return lm.getLogger(name)
}

/**
 * All appenders should register them in their module init() method.
 * The method returns error if the function is called after config intialization sub-system.
 * Parameters:
 *		appenderName - name of the appender, this name can be present in config file
 *		newAppender - function which allows to create the appender instance. It
 *				expects map of named params to their values.
 */
func RegisterAppender(appenderName string, newAppenderFn NewAppenderFn) error {
	return lm.registerAppender(appenderName, newAppenderFn)
}

/**
 * Should be called to shutdown log subsystem properly. It will notify all logContexts and wait
 * while all go routines are over. To call this method could be essential to finalize some 
 * appenders implementations and close them properly
 */
func Shutdown() {
	// TODO: implement it
}