package log4g

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// log4g the appender registration name
const fileAppenderName = "log4g/fileAppender"

// layout - appender setting to specify format of the log event to message transformation
// This param must be provided when new appender is created
const FAParamLayout = "layout"

// fileName - specifies fileName where log message will be written to
// This param must be provided when new appender is created
const FAParamFileName = "fileName"

// append - appender setting if it is set to true the new log messages will be added to
// the end of file appending them to the current content of the file.
// this parameter is OPTIONAL, default value is true
const FAParamFileAppend = "append"

// buffer - appender settings which allows to set the size of log event buffer
// this parameter is OPTIONAL, default value is 100
const FAParamFileBuffer = "buffer"

// maxFileSize - appender settings which limits the size of the file chunk. The size can be specified in
// human readable form like 10Mb or 1000kiB etc. 1024 <= maxFileSize <= maxInt64
// this parameter is OPTIONAL, default value is maxInt64. It is ignored if rotate != maxSize
const FAParamMaxSize = "maxFileSize"

// maxLines - appender settings which limits the number of lines in a file chunk. The number can be specified in
// human readable form like 10Mb or 400kB etc. 1000 <= maxFileSize <= maxInt64
// this parameter is OPTIONAL, default value is maxInt64. It is ignored if rotate != maxLines
const FAParamMaxLines = "maxLines"

// rotate - defines file-chunks rotation policy.
// this parameter is OPTIONAL, default value is none.
const FAParamRotate = "rotate"

// possible values of rotate param
// none: no rotation at all
// size: just rotate if maxFileSize OR maxLines is reached
// daily: rotate every new day (the host time midnight) or maxFileSize OR maxLines is reached
var rotateState = map[string]int{"none": rsNone, "size": rsSize, "daily": rsDaily}

const (
	rsNone = iota
	rsSize
	rsDaily
)

type fileAppenderFactory struct {
}

type fileAppender struct {
	msgChannel     chan string
	controlCh      chan bool
	fileName       string
	file           *os.File
	layoutTemplate LayoutTemplate
	fileAppend     bool
	maxSize        int64
	maxLines       int64
	rotate         int
	stat           stats
}

type stats struct {
	lines         int64
	size          int64
	startTime     time.Time
	lastErrorTime time.Time
}

func init() {
	RegisterAppender(&fileAppenderFactory{})
}

// The factory allows to create an appender instances
func (faf *fileAppenderFactory) Name() string {
	return fileAppenderName
}

func (faf *fileAppenderFactory) NewAppender(params map[string]string) (Appender, error) {
	layout, ok := params[FAParamLayout]
	if !ok || len(layout) == 0 {
		return nil, errors.New("Cannot create file appender: layout should be specified")
	}

	fileName, ok := params[FAParamFileName]
	if !ok || len(fileName) == 0 {
		return nil, errors.New("Cannot create file appender: file should be specified")
	}

	layoutTemplate, err := ParseLayout(layout)
	if err != nil {
		return nil, errors.New("Cannot create file appender, incorrect layout: " + err.Error())
	}

	buffer, err := ParseInt(params[FAParamFileBuffer], 1, 10000, 100)
	if err != nil {
		return nil, errors.New("Invalid " + FAParamFileBuffer + " value: " + err.Error())
	}

	fileAppend, err := ParseBool(params[FAParamFileAppend], true)
	if err != nil {
		return nil, errors.New("Invalid " + FAParamFileAppend + " value: " + err.Error())
	}

	maxFileSize, err := ParseInt(params[FAParamMaxSize], 1000, maxInt64, maxInt64)
	if err != nil {
		return nil, errors.New("Invalid " + FAParamMaxSize + " value: " + err.Error())
	}

	maxLines, err := ParseInt(params[FAParamMaxLines], 1000, maxInt64, maxInt64)
	if err != nil {
		return nil, errors.New("Invalid " + FAParamMaxLines + " value: " + err.Error())
	}

	rotateStr, ok := params[FAParamRotate]
	rotateStr = strings.Trim(rotateStr, " ")
	rState := rsNone
	if ok && len(rotateStr) > 0 {
		rState, ok = rotateState[rotateStr]
		if !ok {
			return nil, errors.New("Unknown rotate state \"" + rotateStr +
				"\", expected \"none\", \"size\", or \"daily \" value")
		}
	}

	app := &fileAppender{}
	app.msgChannel = make(chan string, buffer)
	app.controlCh = make(chan bool, 1)
	app.layoutTemplate = layoutTemplate
	app.fileName = fileName
	app.fileAppend = fileAppend
	app.maxSize = maxFileSize
	app.maxLines = maxLines
	app.rotate = rState

	go func() {
		defer app.close()
		for {
			str, ok := <-app.msgChannel
			if !ok {
				break
			}

			if app.isRotationNeeded() {
				app.rotateFile()
			}
			app.writeMsg(str)
		}
	}()
	return app, nil
}

func (faf *fileAppenderFactory) Shutdown() {
	// do nothing here, appenders should be shut down by log context
}

func (fa *fileAppender) Append(event *LogEvent) (ok bool) {
	ok = false
	defer EndQuietly()
	msg := ToLogMessage(event, fa.layoutTemplate)
	fa.msgChannel <- msg
	ok = true
	return ok
}

func (fa *fileAppender) Shutdown() {
	close(fa.msgChannel)
	<-fa.controlCh
}

func (fa *fileAppender) rotateFile() error {
	fa.archiveCurrent()

	fa.stat.lines = 0
	fa.stat.size = 0
	fa.stat.startTime = time.Now()

	flags := os.O_WRONLY | os.O_CREATE
	if fa.fileAppend {
		flags = os.O_WRONLY | os.O_APPEND | os.O_CREATE
		fa.stat.lines = fa.linesCount()
		if fInfo, err := os.Stat(fa.fileName); err == nil {
			fa.stat.size = fInfo.Size()
		}
	}
	fmt.Println("new file ", fa.fileAppend)

	fd, err := os.OpenFile(fa.fileName, flags, 0660)
	if err != nil {
		panic("File Appender cannot open file " + fa.fileName + " to store logs: " + err.Error())
	}
	fa.file = fd

	return nil
}

func (fa *fileAppender) linesCount() int64 {
	file, err := os.Open(fa.fileName)
	if err != nil {
		return 0
	}
	defer file.Close()

	var result int64 = 0
	for scanner := bufio.NewScanner(file); scanner.Scan(); result++ {
	}

	return result
}

func (fa *fileAppender) archiveCurrent() {
	if fa.file == nil {
		return
	}

	archiveName, _ := filepath.Abs(fa.fileName)
	if fa.rotate == rsDaily {
		archiveName += "." + fa.stat.startTime.Format("2006-01-02")
	}

	dir := filepath.Dir(archiveName)
	baseName := filepath.Base(archiveName)
	fileInfos, err := ioutil.ReadDir(dir)
	id := 1
	for _, fInfo := range fileInfos {
		if fInfo.IsDir() || !strings.HasPrefix(fInfo.Name(), baseName) {
			continue
		}

		idx := strings.LastIndex(fInfo.Name(), ".")
		if idx < 0 {
			continue
		}

		fId, err := strconv.Atoi(fInfo.Name()[idx+1:])
		if err != nil {
			continue
		}

		if fId >= id {
			id = fId + 1
		}
	}
	fa.file.Close()
	fa.file = nil
	archiveName = archiveName + "." + strconv.Itoa(id)
	err = os.Rename(fa.fileName, archiveName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File appender %+v: it is impossible to rename file \"%s\" to \"%s\": %s\n", fa, fa.fileName, archiveName, err)
	}
}

func (fa *fileAppender) isRotationNeeded() bool {
	if fa.file == nil {
		return true
	}

	if fa.stat.size > fa.maxSize || fa.stat.lines > fa.maxLines {
		return true
	}

	if fa.rotate != rsDaily {
		return false
	}

	now := time.Now()
	return fa.stat.startTime.Day() != now.Day() || now.Sub(fa.stat.startTime) > time.Hour*24
}

func (fa *fileAppender) writeMsg(msg string) {
	n, err := fmt.Fprint(fa.file, msg, "\n")

	if err != nil {
		if time.Since(fa.stat.lastErrorTime) > time.Minute {
			fa.stat.lastErrorTime = time.Now()
			fmt.Fprintf(os.Stderr, "File appender %+v: %s\n", fa, err)
		}
		return
	}

	fa.stat.lines++
	fa.stat.size += int64(n)
}

func (fa *fileAppender) close() {
	err := recover()
	if err != nil {
		fmt.Fprintf(os.Stderr, "File appender %+v: %s\n", fa, err)
	}
	if fa.file != nil {
		fa.file.Close()
		fa.file = nil
	}
	fa.controlCh <- true
	close(fa.controlCh)
}
