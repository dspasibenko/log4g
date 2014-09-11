package log4g

import "time"

// Level type represents logging level as an integer in [0..70] range.
// A level with lowest value has higher priority than a level with highest value.
type Level int

// Predefined log levels. Users can define own ones or overwrite the predefined via configuration
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

// Logger interface provides methods for delivering messages to different appenders
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

// LogEvent is DTO, bearing a log message to log (final destination of the message)
type LogEvent struct {
	Level      Level
	Timestamp  time.Time
	LoggerName string
	Payload    interface{}
}

// Appender is an interface for a log endpoint. Different storages can be connected to the library
// implementing the interface
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

// GetLogger returns pointer to the Logger object for specified logger name.
// The function will always return the same pointer for the same logger's name
// regardless of log4g configuration or other settings
func GetLogger(loggerName string) Logger {
	return lm.getLogger(loggerName)
}

func SetLogLevel(loggerName string, level Level) {
	lm.setLogLevel(loggerName, level)
}

// RegisterAppender allows to register an appender implementation in log4g. All appenders should register themself calling the
// function from init() or by calling this function directly.
// The method returns error if another factory has been registered for the same name before
// Parameters:
//		appenderFactory - a factory object which allows to create new instances of the appender type.
func RegisterAppender(appenderFactory AppenderFactory) error {
	return lm.registerAppender(appenderFactory)
}

// ConfigF reads log4g configuration properties from text file, which name is provided in
// configFileName parameter.
func ConfigF(configFileName string) error {
	return lm.setPropsFromFile(configFileName)
}

// Config allows to configure log4g by properties provided in the key:value form
func Config(props map[string]string) error {
	return lm.setNewProperties(props)
}

// Should be called to shutdown log subsystem properly. It will notify all logContexts and wait
// while all go routines that deliver messages to appenders are over. Calling this method could
// be essential to finalize some appenders and release their resources properly
func Shutdown() {
	lm.shutdown()
}
