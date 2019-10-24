package xlog

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type Logger struct {
	mu        sync.RWMutex
	out       Output
	severity  Severity
	verbose   Verbose
	verbosity Verbose
	tm        time.Time
	fields    Fields
}

func New(out Output, severity Severity, verbose Verbose) *Logger {
	if !severity.IsValid() {
		severity = SeverityInfo
	}
	return &Logger{
		out:      out,
		severity: severity,
		verbose:  verbose,
	}
}

func (l *Logger) clone() *Logger {
	l.mu.RLock()
	ln := &Logger{
		out:       l.out,
		severity:  l.severity,
		verbose:   l.verbose,
		verbosity: l.verbosity,
		tm:        l.tm,
		fields:    make(Fields, len(l.fields)),
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
	if l.out != nil && l.severity >= severity && l.verbose >= l.verbosity {
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
		l.out.Log(buf, severity, l.verbosity, tm, l.fields, callers)
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

func (l *Logger) SetOutput(out Output) {
	l.mu.Lock()
	l.out = out
	l.mu.Unlock()
}

func (l *Logger) SetSeverity(severity Severity) {
	l.mu.Lock()
	if !severity.IsValid() {
		severity = SeverityInfo
	}
	l.severity = severity
	l.mu.Unlock()
}

func (l *Logger) SetVerbose(verbose Verbose) {
	l.mu.Lock()
	l.verbose = verbose
	l.mu.Unlock()
}

func (l *Logger) V(verbosity Verbose) *Logger {
	ln := l.clone()
	ln.verbosity = verbosity
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
