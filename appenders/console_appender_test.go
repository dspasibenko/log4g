package appenders

import (
	"github.com/dspasibenko/log4g"
	. "gopkg.in/check.v1"
	"time"
)

type cAppenderSuite struct {
	msg    string
	signal chan bool
}

var _ = Suite(&cAppenderSuite{})

func (cas *cAppenderSuite) Write(p []byte) (n int, err error) {
	cas.msg = string(p)
	cas.signal <- true
	return len(p), nil
}

func (s *cAppenderSuite) TestNewAppender(c *C) {
	a, err := caFactory.NewAppender(map[string]interface{}{})
	c.Assert(a, IsNil)
	c.Assert(err, NotNil)

	a, err = caFactory.NewAppender(map[string]interface{}{"abcd": "1234"})
	c.Assert(a, IsNil)
	c.Assert(err, NotNil)

	a, err = caFactory.NewAppender(map[string]interface{}{"layout": "%c %p"})
	c.Assert(a, NotNil)
	c.Assert(err, IsNil)
}

func (s *cAppenderSuite) TestAppend(c *C) {
	s.signal = make(chan bool, 1)
	caFactory.out = s

	a, _ := caFactory.NewAppender(map[string]interface{}{"layout": "[%d{15:04:05.000}] %p %c: %m"})
	appended := a.Append(&log4g.LogEvent{log4g.FATAL, time.Unix(123456, 0), "a.b.c", "Hello Console!"})
	c.Assert(appended, Equals, true)
	<-s.signal
	c.Assert(s.msg, Equals, "[02:17:36.000] FATAL a.b.c: Hello Console!")

	caFactory.Shutdown()
	appended = a.Append(&log4g.LogEvent{log4g.FATAL, time.Unix(0, 0), "a.b.c", "Never delivered"})
	c.Assert(appended, Equals, false)
}
