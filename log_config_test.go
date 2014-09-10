package log4g

import (
	. "gopkg.in/check.v1"
	//"time"
)

type logConfigSuite struct {
}

var _ = Suite(&logConfigSuite{})

func (s *logConfigSuite) TestGetLevelByName(c *C) {
	lc := newLogConfig()
	lc.applyLevelParams(map[string]string{})
	c.Assert(lc.getLevelByName(lc.levelNames[FATAL]), Equals, FATAL)
	c.Assert(lc.getLevelByName(lc.levelNames[ERROR]), Equals, ERROR)
	c.Assert(lc.getLevelByName(lc.levelNames[WARN]), Equals, WARN)
	c.Assert(lc.getLevelByName(lc.levelNames[INFO]), Equals, INFO)
	c.Assert(lc.getLevelByName(lc.levelNames[DEBUG]), Equals, DEBUG)
	c.Assert(lc.getLevelByName(lc.levelNames[TRACE]), Equals, TRACE)
}

func (s *logConfigSuite) TestMergedParamsWithDefault(c *C) {
	params := mergedParamsWithDefault(map[string]string{"abcd": "efgh"})
	c.Assert(params["abcd"], Equals, "efgh")
	for k, v := range defaultConfigParams {
		c.Assert(params[k], Equals, v)
		c.Assert(len(v), Not(Equals), 0)
	}
}

func (s *logConfigSuite) TestInitIfNeeded(c *C) {
	f := &consoleAppenderFactory{"log4g/appenders/consoleAppender"}
	RegisterAppender(f)
	lc := lm.config
	lc.initIfNeeded()

	c.Assert(len(lc.loggers), Equals, 0)
	c.Assert(lc.logLevels.At(0).(*logLevelSetting).level, Equals, INFO)
	c.Assert(lc.logContexts.At(0).(*logContext).appenders[0], NotNil)
	c.Assert(lc.appenderFactorys[f.Name()], Equals, f)
	c.Assert(len(lc.appenders), Equals, 1)

	lc.logLevels.At(0).(*logLevelSetting).level = DEBUG
	lc.initIfNeeded()
	c.Assert(lc.logLevels.At(0).(*logLevelSetting).level, Equals, DEBUG)
}

// Console Appender mocking stuff
type consoleAppender struct {
}

type consoleAppenderFactory struct {
	name string
}

func (caf *consoleAppenderFactory) Name() string {
	return caf.name
}

func (caf *consoleAppenderFactory) NewAppender(params map[string]string) (Appender, error) {
	return &consoleAppender{}, nil
}

func (caf *consoleAppenderFactory) Shutdown() {
}

// Appender interface implementation
func (cAppender *consoleAppender) Append(event *LogEvent) (ok bool) {
	return true
}

func (cAppender *consoleAppender) Shutdown() {

}
