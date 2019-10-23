package xlog

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type Logger struct {
	mu          sync.RWMutex
	out         LogOutput
	maxSeverity Severity
	maxVerbose  Verbose
	verbose     Verbose
	tm          time.Time
	fields      Fields
}

func New(out LogOutput, maxSeverity Severity, maxVerbose Verbose) *Logger {
	if sv := int(maxSeverity); sv < 0 || sv >= len(severities) {
		maxSeverity = SeverityInfo
	}
	return &Logger{
		out:         out,
		maxSeverity: maxSeverity,
		maxVerbose:  maxVerbose,
	}
}

func (l *Logger) clone() *Logger {
	l.mu.RLock()
	ln := &Logger{
		out:         l.out,
		maxSeverity: l.maxSeverity,
		maxVerbose:  l.maxVerbose,
		verbose:     l.verbose,
		tm:          l.tm,
		fields:      make(Fields, len(l.fields)),
	}
	for key := range l.fields {
		ln.fields[key] = l.fields[key]
	}
	l.mu.RUnlock()
	return ln
}

func (l *Logger) output(severity Severity, message string) {
	if l == nil {
		return
	}
	l.mu.RLock()
	if l.out != nil && l.maxSeverity >= severity && l.maxVerbose >= l.verbose {
		messageLen := len(message)
		buf := make([]byte, 0, messageLen+1)
		buf = append(buf, message...)
		if messageLen == 0 || message[messageLen-1] != '\n' {
			buf = append(buf, '\n')
		}
		tm := l.tm
		if tm.IsZero() {
			tm = time.Now()
		}
		callers := make([]uintptr, 32)
		callers = callers[:runtime.Callers(4, callers)]
		l.out.Log(buf, severity, l.verbose, tm, l.fields, runtime.CallersFrames(callers))
	}
	l.mu.RUnlock()
}

func (l *Logger) log(severity Severity, args ...interface{}) {
	l.output(severity, fmt.Sprint(args...))
}

func (l *Logger) logf(severity Severity, format string, args ...interface{}) {
	l.output(severity, fmt.Sprintf(format, args...))
}

func (l *Logger) logln(severity Severity, args ...interface{}) {
	l.output(severity, fmt.Sprintln(args...))
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log(SeverityFatal, args...)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logf(SeverityFatal, format, args...)
	os.Exit(1)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.logln(SeverityFatal, args...)
	os.Exit(1)
}

func (l *Logger) Error(args ...interface{}) {
	l.log(SeverityError, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(SeverityError, format, args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.logln(SeverityError, args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.log(SeverityWarning, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.logf(SeverityWarning, format, args...)
}

func (l *Logger) Warningln(args ...interface{}) {
	l.logln(SeverityWarning, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log(SeverityInfo, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(SeverityInfo, format, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.logln(SeverityInfo, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(SeverityDebug, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(SeverityDebug, format, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.logln(SeverityDebug, args...)
}

func (l *Logger) SetOutput(out LogOutput) {
	l.mu.Lock()
	l.out = out
	l.mu.Unlock()
}

func (l *Logger) SetMaxSeverity(maxSeverity Severity) {
	l.mu.Lock()
	if sv := int(maxSeverity); sv < 0 || sv >= len(severities) {
		maxSeverity = SeverityInfo
	}
	l.maxSeverity = maxSeverity
	l.mu.Unlock()
}

func (l *Logger) SetMaxVerbose(maxVerbose Verbose) {
	l.mu.Lock()
	l.maxVerbose = maxVerbose
	l.mu.Unlock()
}

func (l *Logger) V(verbose Verbose) *Logger {
	ln := l.clone()
	ln.verbose = verbose
	return ln
}

func (l *Logger) WithTime(tm time.Time) *Logger {
	ln := l.clone()
	ln.tm = tm
	return ln
}

func (l *Logger) WithFields(fields Fields) *Logger {
	ln := l.clone()
	for key, val := range ln.fields {
		ln.fields[key] = val
	}
	return ln
}
