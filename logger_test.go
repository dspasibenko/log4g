package log4g

import (
	. "gopkg.in/check.v1"
)

type loggerSuite struct {
	loggerName string
}

var _ = Suite(&loggerSuite{})

func (s *loggerSuite) TestApplyNewLevelToLoggers(c *C) {
	rootLLS := &logLevelSetting{rootLoggerName, INFO}

	loggers := make(map[string]*logger)
	loggers["a"] = &logger{"a", rootLLS, INFO}
	loggers["a.b"] = &logger{"a.b", rootLLS, INFO}
	loggers["a.b.c"] = &logger{"a.b.c", rootLLS, INFO}
	loggers["a.b.c.d"] = &logger{"a.b.c.d", rootLLS, INFO}

	applyNewLevelToLoggers(&logLevelSetting{"a.b", DEBUG}, loggers)
	c.Assert(loggers["a"].logLevel, Equals, INFO)
	c.Assert(loggers["a.b"].logLevel, Equals, DEBUG)
	c.Assert(loggers["a.b.c"].logLevel, Equals, DEBUG)
	c.Assert(loggers["a.b.c.d"].logLevel, Equals, DEBUG)

	applyNewLevelToLoggers(&logLevelSetting{"a.b.c", TRACE}, loggers)
	applyNewLevelToLoggers(&logLevelSetting{"a.b", ERROR}, loggers)
	c.Assert(loggers["a"].logLevel, Equals, INFO)
	c.Assert(loggers["a.b"].logLevel, Equals, ERROR)
	c.Assert(loggers["a.b.c"].logLevel, Equals, TRACE)
	c.Assert(loggers["a.b.c.d"].logLevel, Equals, TRACE)
}
