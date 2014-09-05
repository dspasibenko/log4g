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

func (nus *nameUtilsSuite) name() string {
	return nus.loggerName
}

func (nus *nameUtilsSuite) Compare(other collections.Comparator) int {
	return compare(nus, other.(*nameUtilsSuite))
}
