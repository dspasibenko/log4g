package log4g

import (
	"github.com/dspasibenko/log4g/collections"
	"regexp"
	"strings"
)

type logNameProvider interface {
	name() string
}

// Config params
const (
	cfgAppender     = "appender."
	cfgAppenderType = "type"
	cfgLevel        = "level."
	cfgContext      = "context."
)

func compare(n1, n2 logNameProvider) int {
	switch {
	case n1.name() == n2.name():
		return 0
	case n1.name() < n2.name():
		return -1
	}
	return 1
}

// loggerName cannot start/end from spaces and dots
func normalizeLogName(name string) string {
	return strings.ToLower(strings.Trim(name, ". "))
}

/**
 * Checks whether the checkedName is ancestor for the loggerName or not
 * The name checkedName is ancestor for the loggerName if:
 *	- checkedName == loggerName
 *  - loggerName == checkedName.<some name here>
 * 	- checkedName == rootLoggerName
 */
func ancestor(checkedName, loggerName string) bool {
	if checkedName == loggerName || checkedName == rootLoggerName {
		return true
	}

	lenc := len(checkedName)
	lenl := len(loggerName)
	if strings.HasPrefix(loggerName, checkedName) && lenl > lenc && loggerName[lenc] == '.' {
		return true
	}
	return false
}

func getNearestAncestor(comparator collections.Comparator, names *collections.SortedSlice) logNameProvider {
	if names.Len() == 0 {
		return nil
	}
	nProvider := comparator.(logNameProvider)
	for idx := Min(names.Len()-1, names.GetInsertPos(nProvider.(collections.Comparator))); idx >= 0; idx-- {
		candidate := names.At(idx).(logNameProvider)
		if ancestor(candidate.name(), nProvider.name()) {
			return candidate
		}
	}
	return nil
}

// parses param and expect the appender name in the form: appender.<appenderName>.<appenderParam>
func getAppenderName(param string) (string, bool) {
	p := cfgAppender
	if !strings.HasPrefix(param, p) {
		return "", false
	}

	start := len(p)
	end := strings.LastIndex(param, ".")
	if start >= end-1 {
		return "", false
	}

	appenderName := param[start:end]
	if !isCorrectAppenderName(appenderName) {
		return "", false
	}

	return appenderName, true
}

// invariant: param has "appender." prefix
func getAppenderParam(param string) string {
	end := strings.LastIndex(param, ".")
	if end == len(param)-1 {
		return ""
	}
	return param[end+1:]
}

// parses param and expects the context logger name in the form: context.<loggerName>.<contextParam>
func getContextLoggerName(param string) (string, bool) {
	c := cfgContext
	if !strings.HasPrefix(param, c) {
		return "", false
	}

	start := len(c)
	end := strings.LastIndex(param, ".")
	if start == end+1 {
		return "", true
	}

	loggerName := param[start:end]
	if !isCorrectLoggerName(loggerName) {
		return "", false
	}

	return loggerName, true
}

// invariant: param has "context." prefix
func getContextParam(param string) string {
	end := strings.LastIndex(param, ".")
	if end == len(param)-1 {
		return ""
	}
	return param[end+1:]
}

func parseAppendersParams(params map[string]string) map[string]map[string]string {
	// collect settings for all appenders from config
	apps := make(map[string]map[string]string)
	for k, v := range params {
		appName, ok := getAppenderName(k)
		if !ok {
			continue
		}
		appParam := getAppenderParam(k)

		m, ok := apps[appName]
		if !ok {
			m = make(map[string]string)
			apps[appName] = m
		}
		m[appParam] = v
	}
	return apps
}

func parseContextParams(params map[string]string) map[string]map[string]string {
	ctxts := make(map[string]map[string]string)
	for k, v := range params {
		logName, ok := getContextLoggerName(k)
		if !ok {
			continue
		}
		ctxParam := getContextParam(k)

		m, ok := ctxts[logName]
		if !ok {
			m = make(map[string]string)
			ctxts[logName] = m
		}
		m[ctxParam] = v
	}
	return ctxts
}

func isCorrectAppenderName(appenderName string) bool {
	matched, err := regexp.MatchString("^[A-Za-z][A-Za-z0-9.]+$", appenderName)
	if !matched || err != nil {
		return false
	}
	return true
}

func isCorrectLoggerName(loggerName string) bool {
	matched, err := regexp.MatchString("^[A-Za-z]+([A-Za-z0-9.]*[A-Za-z0-9]+)*$", loggerName)
	if !matched || err != nil {
		return false
	}
	return true
}

// Utility methods
func Min(a, b int) int {
	if a < b {
		return a
	} else if b < a {
		return b
	}
	return a
}

func EndQuietly() {
	recover()
}
