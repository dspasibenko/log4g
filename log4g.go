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
}

type Timestamp uint64

type LogEvent struct {
	level     Level
	timestamp Timestamp
	logger    string
	message   string
}

type Appender interface {
	Append(record *LogEvent)
}

/**
 * Provides pointer to the logger with specified name.
 * name can have 'dot' separated form.
 */
func GetLogger(name string) Logger {
	return lm.getLogger(name)
}
