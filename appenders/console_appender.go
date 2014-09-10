package appenders

import (
	"errors"
	"fmt"
	"github.com/dspasibenko/log4g"
	"io"
	"os"
)

const consoleAppenderName = "log4g/appenders/consoleAppender"

// Parameters accepted by the appender
const CAParamLayout = "layout"

type consoleAppender struct {
	layoutTemplate LayoutTemplate
}

type consoleAppenderFactory struct {
	msgChannel chan string
	out        io.Writer
}

var caFactory *consoleAppenderFactory

func InitConsoleAppender() {
	caFactory = &consoleAppenderFactory{make(chan string, 1000), os.Stdout}
	err := log4g.RegisterAppender(caFactory)
	if err != nil {
		close(caFactory.msgChannel)
		fmt.Println("It is impossible to register console appender: ", err)
		panic(err)
	}
	go func() {
		for {
			str, ok := <-caFactory.msgChannel
			if !ok {
				break
			}
			fmt.Fprint(caFactory.out, str)
		}
	}()
}

func (*consoleAppenderFactory) Name() string {
	return consoleAppenderName
}

func (caf *consoleAppenderFactory) NewAppender(params map[string]string) (log4g.Appender, error) {
	layout, ok := params[CAParamLayout]
	if !ok || len(layout) == 0 {
		return nil, errors.New("Cannot create console appender without specified layout")
	}

	layoutTemplate, err := ParseLayout(layout)
	if err != nil {
		return nil, err
	}

	return &consoleAppender{layoutTemplate}, nil
}

func (caf *consoleAppenderFactory) Shutdown() {
	close(caf.msgChannel)
}

// Appender interface implementation
func (cAppender *consoleAppender) Append(event *log4g.LogEvent) (ok bool) {
	ok = false
	defer log4g.EndQuietly()
	msg := ToLogMessage(event, cAppender.layoutTemplate)
	caFactory.msgChannel <- msg
	ok = true
	return ok
}

func (cAppender *consoleAppender) Shutdown() {

}
