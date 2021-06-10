package xlog

import (
	"time"

	"github.com/goinsane/erf"
)

// Log carries the log.
type Log struct {
	Message     []byte
	Severity    Severity
	Verbosity   Verbose
	Time        time.Time
	Fields      Fields
	StackCaller erf.StackCaller
	StackTrace  *erf.StackTrace
}

// Duplicate duplicates the Log.
func (l *Log) Duplicate() *Log {
	if l == nil {
		return nil
	}
	l2 := &Log{
		Message:     make([]byte, len(l.Message)),
		Severity:    l.Severity,
		Verbosity:   l.Verbosity,
		Time:        l.Time,
		Fields:      l.Fields.Duplicate(),
		StackCaller: l.StackCaller,
		StackTrace:  l.StackTrace.Duplicate(),
	}
	copy(l2.Message, l.Message)
	return l2
}
