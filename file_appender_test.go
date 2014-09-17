package log4g

import (
	. "gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type faConfigSuite struct {
}

var _ = Suite(&faConfigSuite{})

func (s *faConfigSuite) TestFileAppenderName(c *C) {
	c.Assert(faFactory.Name(), Equals, "log4g/fileAppender")
}

func (s *faConfigSuite) TestNewAppender(c *C) {
	app, err := faFactory.NewAppender(nil)
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %ee"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %ee", "fileName": "fn"}) //bad layout
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "-1"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000", "append": "abc"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000", "append": "true", "maxFileSize": "10"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000", "append": "true",
		"maxFileSize": "2K", "maxLines": "10"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000", "append": "true",
		"maxFileSize": "2K", "maxLines": "2G", "rotate": "daily2"})
	c.Assert(app, IsNil)
	c.Assert(err, NotNil)

	app, err = faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000", "append": "true",
		"maxFileSize": "2K", "maxLines": "2M", "rotate": "size"})
	c.Assert(app, NotNil)
	c.Assert(err, IsNil)
	c.Assert(app.(*fileAppender).rotate, Equals, rsSize)
	app.Shutdown()
}

func (s *faConfigSuite) TestShutdown(c *C) {
	app, err := faFactory.NewAppender(map[string]string{"layout": " %p", "fileName": "fn", "buffer": "1000",
		"maxFileSize": "2K", "maxLines": "2M", "rotate": "daily"})
	c.Assert(app, NotNil)
	c.Assert(err, IsNil)
	fa := app.(*fileAppender)
	c.Assert(fa.file, IsNil)
	c.Assert(fa.fileAppend, Equals, true)
	c.Assert(fa.fileName, Equals, "fn")
	c.Assert(fa.maxSize, Equals, int64(2000))
	c.Assert(fa.maxLines, Equals, int64(2000000))
	c.Assert(fa.rotate, Equals, rsDaily)
	ok := false
	select {
	case _, ok = <-fa.msgChannel:
		ok = true
		break
	default:
	}
	c.Assert(ok, Equals, false)
	app.Shutdown()

	_, ok = <-fa.msgChannel
	c.Assert(ok, Equals, false)
}

func (s *faConfigSuite) TestAppendAndCountLines(c *C) {
	defer os.Remove("fn")
	fa := writeLogs(c, map[string]string{"layout": "%p", "fileName": "fn", "buffer": "1000",
		"maxFileSize": "2Gib", "maxLines": "2M", "rotate": "daily"}, 10000)
	c.Check(fa.linesCount(), Equals, int64(10000))
}

func (s *faConfigSuite) TestAppendToExistingOne(c *C) {
	defer removeFiles("____test____log___file")
	params := map[string]string{"layout": "%p", "fileName": "____test____log___file", "buffer": "1000",
		"maxFileSize": "2Gib", "maxLines": "2M", "rotate": "none", "append": "true"}
	writeLogs(c, params, 10000)
	fa := writeLogs(c, params, 10000)
	c.Check(fa.linesCount(), Equals, int64(20000))
	params["append"] = "false"
	fa = writeLogs(c, params, 550)
	c.Check(fa.linesCount(), Equals, int64(550))
}

func removeFiles(prefix string) {
	archiveName, _ := filepath.Abs(prefix)
	dir := filepath.Dir(archiveName)
	baseName := filepath.Base(archiveName)
	fileInfos, _ := ioutil.ReadDir(dir)
	for _, fInfo := range fileInfos {
		if fInfo.IsDir() || !strings.HasPrefix(fInfo.Name(), baseName) {
			continue
		}
		os.Remove(fInfo.Name())
	}
}

func writeLogs(c *C, params map[string]string, count int) *fileAppender {
	app, _ := faFactory.NewAppender(params)
	c.Assert(app, NotNil)
	fa := app.(*fileAppender)
	for idx := 0; idx < count; idx++ {
		app.Append(&LogEvent{INFO, time.Now(), "abc", "def"})
	}
	app.Shutdown()
	return fa
}
