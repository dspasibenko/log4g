package log4g

import (
	"github.com/dspasibenko/log4g/collections"
	. "gopkg.in/check.v1"
)

type nameUtilsSuite struct {
	loggerName string
}

var _ = Suite(&nameUtilsSuite{})

func (s *nameUtilsSuite) TestAncestor(c *C) {
	c.Assert(ancestor("", ""), Equals, true)
	c.Assert(ancestor("a", "a"), Equals, true)
	c.Assert(ancestor("a.b", "a.b.c"), Equals, true)
	c.Assert(ancestor("a.b", "a.b.cd.e"), Equals, true)
	c.Assert(ancestor("a.b", "a.c.c"), Equals, false)
}

func (s *nameUtilsSuite) TestGetSetLogLevel(c *C) {
	ss, _ := collections.NewSortedSlice(2)
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a"}, ss), IsNil)

	ss.Add(&nameUtilsSuite{"b"})
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a"}, ss), IsNil)

	ss.Add(&nameUtilsSuite{""})
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a"}, ss).(*nameUtilsSuite).loggerName, Equals, "")

	ss.Add(&nameUtilsSuite{"a.b.c"})
	ss.Add(&nameUtilsSuite{"a.b"})
	ss.Add(&nameUtilsSuite{"a"})
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a.b.d"}, ss).(*nameUtilsSuite).loggerName, Equals, "a.b")
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a.b.c"}, ss).(*nameUtilsSuite).loggerName, Equals, "a.b.c")
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a.b.c.d"}, ss).(*nameUtilsSuite).loggerName, Equals, "a.b.c")
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a.b.c.d"}, ss).(*nameUtilsSuite).loggerName, Equals, "a.b.c")
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a.bc.d"}, ss).(*nameUtilsSuite).loggerName, Equals, "a")
}

func (s *nameUtilsSuite) TestGetAppenderName(c *C) {
	a, ok := getAppenderName("abc")
	c.Assert(ok, Equals, false)

	a, ok = getAppenderName("abc.asd.ab")
	c.Assert(ok, Equals, false)

	a, ok = getAppenderName("appender.test")
	c.Assert(ok, Equals, false)

	a, ok = getAppenderName("appender..test")
	c.Assert(ok, Equals, false)

	a, ok = getAppenderName("appender...test")
	c.Assert(ok, Equals, false)

	a, ok = getAppenderName("appender.test.test")
	c.Assert(ok, Equals, true)
	c.Assert(a, Equals, "test")

	a, ok = getAppenderName("appender.3test.test")
	c.Assert(ok, Equals, false)

	a, ok = getAppenderName("appender.t345est.test")
	c.Assert(ok, Equals, true)
	c.Assert(a, Equals, "t345est")
}

func (s *nameUtilsSuite) TestGetAppenderParam(c *C) {
	c.Assert(getAppenderParam("appender."), Equals, "")
	c.Assert(getAppenderParam("appender.ROOT.level"), Equals, "level")
}

func (s *nameUtilsSuite) TestGetContextLoggerName(c *C) {
	ctx, ok := getContextLoggerName("abc")
	c.Assert(ok, Equals, false)

	ctx, ok = getContextLoggerName("abc.asd.ab")
	c.Assert(ok, Equals, false)

	ctx, ok = getContextLoggerName("context.test")
	c.Assert(ok, Equals, true)
	c.Assert(ctx, Equals, "")

	ctx, ok = getContextLoggerName("context..test")
	c.Assert(ok, Equals, false)

	ctx, ok = getContextLoggerName("context...test")
	c.Assert(ok, Equals, false)

	ctx, ok = getContextLoggerName("context.test.test")
	c.Assert(ok, Equals, true)
	c.Assert(ctx, Equals, "test")

	ctx, ok = getContextLoggerName("context.a.b.c.test")
	c.Assert(ok, Equals, true)
	c.Assert(ctx, Equals, "a.b.c")
}

func (s *nameUtilsSuite) TestParseAppendersParams(c *C) {
	params := parseAppendersParams(map[string]string{
		"appender.ROOT.type": "123",
		"abc":                "def",
		"appender.app.type":  "345",
		"appender.ROOT.ttt":  "qqq",
	})
	c.Assert(params["ROOT"]["type"], Equals, "123")
	c.Assert(params["ROOT"]["ttt"], Equals, "qqq")
	c.Assert(params["app"]["type"], Equals, "345")
	c.Assert(params["app"]["ttt"], Equals, "")
	c.Assert(params["abc"], IsNil)
}

func (s *nameUtilsSuite) TestParseContextParams(c *C) {
	params := parseContextParams(map[string]string{
		"context.ROOT.type": "123",
		"abc":               "def",
		"context.app.type":  "345",
		"context.ROOT.ttt":  "qqq",
	})
	c.Assert(params["ROOT"]["type"], Equals, "123")
	c.Assert(params["ROOT"]["ttt"], Equals, "qqq")
	c.Assert(params["app"]["type"], Equals, "345")
	c.Assert(params["app"]["ttt"], Equals, "")
	c.Assert(params["abc"], IsNil)
}

func (s *nameUtilsSuite) TestGetContextParam(c *C) {
	c.Assert(getAppenderParam("context."), Equals, "")
	c.Assert(getAppenderParam("context.param"), Equals, "param")
	c.Assert(getAppenderParam("context.a.b.c.param"), Equals, "param")
}

func (s *nameUtilsSuite) TestCorrectAppenderName(c *C) {
	c.Assert(isCorrectAppenderName(""), Equals, false)
	c.Assert(isCorrectAppenderName("AbcL"), Equals, true)
	c.Assert(isCorrectAppenderName("abC1"), Equals, true)
	c.Assert(isCorrectAppenderName("2abC1"), Equals, false)
	c.Assert(isCorrectAppenderName("ad,CD"), Equals, false)
}

func (s *nameUtilsSuite) TestCorrectLoggerName(c *C) {
	c.Assert(isCorrectLoggerName(""), Equals, false)
	c.Assert(isCorrectLoggerName("a"), Equals, true)
	c.Assert(isCorrectLoggerName("a1"), Equals, true)
	c.Assert(isCorrectLoggerName("1a"), Equals, false)
	c.Assert(isCorrectLoggerName("a.b.c"), Equals, true)
	c.Assert(isCorrectLoggerName(".a"), Equals, false)
	c.Assert(isCorrectLoggerName("a."), Equals, false)
}

func (nus *nameUtilsSuite) name() string {
	return nus.loggerName
}

func (nus *nameUtilsSuite) Compare(other collections.Comparator) int {
	return compare(nus, other.(*nameUtilsSuite))
}
