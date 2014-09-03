package log4g

import "fmt"

type logger struct {
}

func (l *logger) Fatal(args ...interface{}) {
	fmt.Println(args)
}

func (l *logger) Error(args ...interface{}) {

}

func (l *logger) Warn(args ...interface{}) {

}

func (l *logger) Info(args ...interface{}) {

}

func (l *logger) Debug(args ...interface{}) {

}

func (l *logger) Trace(args ...interface{}) {

}

func (l *logger) Logf(level Level, fstr string, args ...interface{}) {

}
