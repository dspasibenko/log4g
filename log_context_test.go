package log4g

import (
	"github.com/dspasibenko/log4g/collections"
	. "gopkg.in/check.v1"
	"time"
)

type logContextSuite struct {
	logEvent *LogEvent
}

var _ = Suite(&logContextSuite{})

func (s *logContextSuite) TestNewLogContext(c *C) {
	lc, err := newLogContext("abc", nil, 10)
	c.Assert(lc, IsNil)
	c.Assert(err, NotNil)

	appenders := make([]Appender, 0, 10)
	lc, err = newLogContext("abc", appenders, 10)
	c.Assert(lc, IsNil)
	c.Assert(err, NotNil)

	appenders = append(appenders, s)
	c.Assert(len(appenders), Equals, 1)
	lc, err = newLogContext("abc", appenders, 0)
	c.Assert(lc, IsNil)
	c.Assert(err, NotNil)
}

func (s *logContextSuite) TestLogContextWorkflow(c *C) {
	appenders := make([]Appender, 1, 10)
	appenders[0] = s
	lc, err := newLogContext("abc", appenders, 1)
	c.Assert(lc, NotNil)
	c.Assert(err, IsNil)

	le := new(LogEvent)
	lc.log(le)
	for i := 0; i < 10; i++ {
		if s.logEvent == le {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	c.Assert(s.logEvent, Equals, le)
	lc.shutdown()

	_, ok := <-lc.controlCh
	c.Assert(ok, Equals, false)
}

func (s *logContextSuite) TestGetLogLevelContext(c *C) {
	ss, _ := collections.NewSortedSlice(2)
	c.Assert(getLogLevelContext("a", ss), IsNil)

	appenders := make([]Appender, 1, 10)
	appenders[0] = s
	lc, _ := newLogContext("b", appenders, 1)
	ss.Add(lc)
	c.Assert(getLogLevelContext("a", ss), IsNil)
	c.Assert(getLogLevelContext("b", ss), Equals, lc)

	lc, _ = newLogContext("", appenders, 1)
	ss.Add(lc)
	c.Assert(getLogLevelContext("a", ss).loggerName, Equals, "")
	c.Assert(getLogLevelContext("b", ss).loggerName, Equals, "b")
}

func (lcs *logContextSuite) Append(logEvent *LogEvent) {
	lcs.logEvent = logEvent
}
