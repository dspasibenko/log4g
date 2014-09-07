package log4g

import (
	"errors"
	"github.com/dspasibenko/log4g/collections"
	"strconv"
)

type logContext struct {
	loggerName string
	appenders  []Appender
	eventsCh   chan *LogEvent
	controlCh  chan bool
}

func newLogContext(loggerName string, appenders []Appender, bufSize int) (*logContext, error) {
	if bufSize <= 0 {
		return nil, errors.New("Cannot create channel with non-positive size=" + strconv.Itoa(bufSize))
	}

	if appenders == nil || len(appenders) == 0 {
		return nil, errors.New("At least one appender should be in Log Context")
	}

	eventsCh := make(chan *LogEvent, bufSize)
	controlCh := make(chan bool, 1)
	lc := &logContext{loggerName, appenders, eventsCh, controlCh}

	go func() {
		defer onStop(controlCh)
		for {
			le, ok := <-eventsCh
			if !ok {
				break
			}
			lc.onEvent(le)
		}
	}()
	return lc, nil
}

// Processing go routine finalizer
func onStop(controlCh chan bool) {
	controlCh <- true
	close(controlCh)
}

func endQuietly() {
	recover()
}

func getLogLevelContext(loggerName string, logContexts *collections.SortedSlice) *logContext {
	lProvider := getNearestAncestor(&logContext{loggerName: loggerName}, logContexts)
	if lProvider == nil {
		return nil
	}
	return lProvider.(*logContext)
}

/**
 * Send the logEvent to all the logContext appenders.
 * Returns true if the logEvent was sent and false if the context is shut down
 */
func (lc *logContext) log(le *LogEvent) (result bool) {
	// Channel can be already closed, so end quietly
	result = false
	defer endQuietly()
	lc.eventsCh <- le
	return true
}

// Called from processing go routine
func (lc *logContext) onEvent(le *LogEvent) {
	for _, a := range lc.appenders {
		a.Append(le)
	}
}

func (lc *logContext) shutdown() {
	close(lc.eventsCh)
	<-lc.controlCh
}

// logNameProvider implementation
func (lc *logContext) name() string {
	return lc.loggerName
}

// Comparator implementation
func (lc *logContext) Compare(other collections.Comparator) int {
	return compare(lc, other.(*logContext))
}
