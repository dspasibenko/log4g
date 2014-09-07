package appenders

import (
	"fmt"
	"github.com/dspasibenko/log4g"
	"io"
	"os"
)

type consoleAppender struct {
}

var stdout io.Writer = os.Stdout

func init() {
	err := log4g.RegisterAppender("log4g/appenders/consoleAppender", newConsoleAppender)
	if err != nil {
		fmt.Println("It is impossible to register console appender: ", err)
		panic(err)
	}
}

// This function will be registered in log4g so it can create this type of appender
func newConsoleAppender(params map[string]interface{}) log4g.Appender {
	return nil
}

// Appender interface implementation
func (console *consoleAppender) Append(event *log4g.LogEvent) {

}
