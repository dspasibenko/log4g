package log4g

import (
	"github.com/dspasibenko/log4g/collections"
	"strings"
)

type logNameProvider interface {
	name() string
}

//TODO: add loggerName checker

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
