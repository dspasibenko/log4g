package log4g

import (
	"github.com/dspasibenko/log4g/collections"
	. "gopkg.in/check.v1"
	"time"
)

type logContextSuite struct {
	logEvents []*LogEvent
	hasSleep  bool
}

var _ = Suite(&logContextSuite{})

func (s *logContextSuite) TestNewLogContext(c *C) {
	lc, err := newLogContext("abc", nil, true, true, 10)
	c.Assert(lc, IsNil)
	c.Assert(err, NotNil)

	appenders := make([]Appender, 0, 10)
	lc, err = newLogContext("abc", appenders, true, true, 10)
	c.Assert(lc, IsNil)
	c.Assert(err, NotNil)

	appenders = append(appenders, s)
	c.Assert(len(appenders), Equals, 1)
	lc, err = newLogContext("abc", appenders, true, true, 0)
	c.Assert(lc, IsNil)
	c.Assert(err, NotNil)
}

func (s *logContextSuite) TestLogContextWorkflow(c *C) {
	appenders := make([]Appender, 1, 10)
	appenders[0] = s
	s.logEvents = make([]*LogEvent, 0, 10)
	lc, err := newLogContext("abc", appenders, true, true, 1)
	c.Assert(lc, NotNil)
	c.Assert(err, IsNil)

	le := new(LogEvent)
	lc.log(le)
	for i := 0; i < 10; i++ {
		if len(s.logEvents) > 0 && s.logEvents[0] == le {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	c.Assert(s.logEvents[0], Equals, le)
	lc.shutdown()

	_, ok := <-lc.controlCh
	c.Assert(ok, Equals, false)
}

func (s *logContextSuite) TestNonBlockingLogContext(c *C) {
	appenders := make([]Appender, 1, 10)
	appenders[0] = s
	s.logEvents = make([]*LogEvent, 0, 100)
	lc, err := newLogContext("abc", appenders, true, false, 1)
	c.Assert(lc, NotNil)
	c.Assert(err, IsNil)

	le := new(LogEvent)
	for i := 0; i < cap(s.logEvents); i++ {
		lc.log(le)
	}

	time.Sleep(time.Millisecond * 50) // give the appender time to read it
	c.Assert(len(s.logEvents), Not(Equals), 0)
	c.Assert(len(s.logEvents), Not(Equals), 100) // non-blocking should lost
	c.Assert(s.logEvents[0], Equals, le)
	lc.shutdown()

	_, ok := <-lc.controlCh
	c.Assert(ok, Equals, false)
}

func (s *logContextSuite) TestGetLogLevelContext(c *C) {
	ss, _ := collections.NewSortedSlice(2)
	c.Assert(getLogLevelContext("a", ss), IsNil)

	appenders := make([]Appender, 1, 10)
	appenders[0] = s
	lc, _ := newLogContext("b", appenders, true, true, 1)
	ss.Add(lc)
	c.Assert(getLogLevelContext("a", ss), IsNil)
	c.Assert(getLogLevelContext("b", ss), Equals, lc)

	lc, _ = newLogContext("", appenders, true, true, 1)
	ss.Add(lc)
	c.Assert(getLogLevelContext("a", ss).loggerName, Equals, "")
	c.Assert(getLogLevelContext("b", ss).loggerName, Equals, "b")
}

func (lcs *logContextSuite) Append(logEvent *LogEvent) bool {
	lcs.logEvents = append(lcs.logEvents, logEvent)
	if lcs.hasSleep {
		time.Sleep(time.Millisecond)
	}
	return true
}

func (lcs *logContextSuite) Shutdown() {
}
