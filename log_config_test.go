package log4g

import (
	. "gopkg.in/check.v1"
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
	lc := newLogConfig()
	c.Assert(lc.registerAppender(&testAppenderFactory{consoleAppenderName}), IsNil)
	lc.initIfNeeded()

	c.Assert(len(lc.loggers), Equals, 0)
	c.Assert(lc.logLevels.At(0).(*logLevelSetting).level, Equals, INFO)
	c.Assert(lc.logContexts.At(0).(*logContext).appenders[0], NotNil)
	c.Assert(lc.appenderFactorys[consoleAppenderName], NotNil)
	c.Assert(len(lc.appenders), Equals, 1)

	lc.logLevels.At(0).(*logLevelSetting).level = DEBUG
	lc.initIfNeeded()
	c.Assert(lc.logLevels.At(0).(*logLevelSetting).level, Equals, DEBUG)
}

func (s *logConfigSuite) TestGetAppendersFromList(c *C) {
	lc := newLogConfig()

	lc.appenders["a"] = &testAppender{"a"}
	lc.appenders["b"] = &testAppender{"b"}
	lc.appenders["c"] = &testAppender{"c"}

	ok := checkPanic(func() { lc.getAppendersFromList("e") })
	c.Assert(ok, Equals, true)

	apps := lc.getAppendersFromList("  ")
	c.Assert(len(apps), Equals, 0)

	apps = lc.getAppendersFromList("b")
	c.Assert(len(apps), Equals, 1)
	c.Assert(apps[0], NotNil)

	apps = lc.getAppendersFromList(" c, a")
	c.Assert(len(apps), Equals, 2)
	c.Assert(apps[0].(*testAppender).name, Equals, "c")
	c.Assert(apps[1].(*testAppender).name, Equals, "a")
}

func (s *logConfigSuite) TestRegisterAppender(c *C) {
	lc := newLogConfig()

	c.Assert(lc.registerAppender(&testAppenderFactory{"a"}), IsNil)
	c.Assert(lc.registerAppender(&testAppenderFactory{"b"}), IsNil)
	c.Assert(lc.registerAppender(&testAppenderFactory{"a"}), NotNil)
	c.Assert(len(lc.appenderFactorys), Equals, 2)
}

func (s *logConfigSuite) TestSetLogLevel(c *C) {
	lc := newLogConfig()
	c.Assert(lc.registerAppender(&testAppenderFactory{consoleAppenderName}), IsNil)
	lc.initIfNeeded()

	l := lc.getLogger("a")

	lc.setLogLevel(DEBUG, "a.b")
	c.Assert(lc.getLogger("a.b.c").(*logger).logLevel, Equals, DEBUG)
	c.Assert(l.(*logger).logLevel, Equals, INFO)

	lc.setLogLevel(WARN, "a")
	c.Assert(l.(*logger).logLevel, Equals, WARN)
	c.Assert(lc.getLogger("a.b.c").(*logger).logLevel, Equals, DEBUG)
}

func (s *logConfigSuite) TestGetLogger(c *C) {
	lc := newLogConfig()
	c.Assert(lc.registerAppender(&testAppenderFactory{consoleAppenderName}), IsNil)
	lc.initIfNeeded()

	l := lc.getLogger("a")

	c.Assert(l, Not(Equals), lc.getLogger("b"))
	c.Assert(l, Equals, lc.getLogger("a"))
	c.Assert(l, Not(Equals), lc.getLogger("A"))
}

func (s *logConfigSuite) TestCreateLoggers(c *C) {
	lc := newLogConfig()
	c.Assert(lc.registerAppender(&testAppenderFactory{consoleAppenderName}), IsNil)
	lc.initIfNeeded()

	params := map[string]string{
		"logger.a.b.c.level": "TRACE",
		"logger.b.c.d.level": "DEBUG",
	}
	lc.createLoggers(params)
	c.Assert(lc.getLogger("a.b.c").(*logger).logLevel, Equals, TRACE)
	c.Assert(lc.getLogger("b.c.d").(*logger).logLevel, Equals, DEBUG)

	pnc := checkPanic(
		func() {
			lc.createLoggers(map[string]string{
				"logger.a.b.c.level": "ABC",
			})
		})
	c.Assert(pnc, Equals, true)
}

type testAppender struct {
	name string
}

type testAppenderFactory struct {
	name string
}

func (taf *testAppenderFactory) Name() string {
	return taf.name
}

func (caf *testAppenderFactory) NewAppender(params map[string]string) (Appender, error) {
	return &testAppender{caf.name}, nil
}

func (caf *testAppenderFactory) Shutdown() {

}

// Appender interface implementation
func (tAppender *testAppender) Append(event *LogEvent) (ok bool) {
	return true
}

func (cAppender *testAppender) Shutdown() {
	// Nothing should be done for the console appender
}
