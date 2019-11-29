package xlog

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

// Logger provides a logger for leveled and structured logging.
type Logger struct {
	mu                 sync.RWMutex
	out                Output
	severity           Severity
	verbose            Verbose
	verbosity          Verbose
	printSeverity      Severity
	stackTraceSeverity Severity
	tm                 time.Time
	fields             Fields
}

// New creates a new Logger. If severity is invalid, it sets SeverityInfo.
func New(out Output, severity Severity, verbose Verbose) *Logger {
	if !severity.IsValid() {
		severity = SeverityInfo
	}
	return &Logger{
		out:                out,
		severity:           severity,
		verbose:            verbose,
		verbosity:          0,
		printSeverity:      SeverityInfo,
		stackTraceSeverity: SeverityNone,
	}
}

func (l *Logger) clone() *Logger {
	l.mu.RLock()
	ln := &Logger{
		out:                l.out,
		severity:           l.severity,
		verbose:            l.verbose,
		verbosity:          l.verbosity,
		printSeverity:      l.printSeverity,
		stackTraceSeverity: l.stackTraceSeverity,
		tm:                 l.tm,
		fields:             l.fields.Clone(),
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
		callers := Callers(nil)
		if l.stackTraceSeverity >= severity {
			callers = make(Callers, 32)
			callers = callers[:runtime.Callers(4, callers)]
		}
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

// Fatal logs to the FATAL severity logs, then calls os.Exit(1).
func (l *Logger) Fatal(args ...interface{}) {
	l.log(SeverityFatal, args...)
	os.Exit(1)
}

// Fatalf logs to the FATAL severity logs, then calls os.Exit(1).
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logf(SeverityFatal, format, args...)
	os.Exit(1)
}

// Fatalln logs to the FATAL severity logs, then calls os.Exit(1).
func (l *Logger) Fatalln(args ...interface{}) {
	l.logln(SeverityFatal, args...)
	os.Exit(1)
}

// Error logs to the ERROR severity logs.
func (l *Logger) Error(args ...interface{}) {
	l.log(SeverityError, args...)
}

// Errorf logs to the ERROR severity logs.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(SeverityError, format, args...)
}

// Errorln logs to the ERROR severity logs.
func (l *Logger) Errorln(args ...interface{}) {
	l.logln(SeverityError, args...)
}

// Warning logs to the WARNING severity logs.
func (l *Logger) Warning(args ...interface{}) {
	l.log(SeverityWarning, args...)
}

// Warningf logs to the WARNING severity logs.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.logf(SeverityWarning, format, args...)
}

// Warningln logs to the WARNING severity logs.
func (l *Logger) Warningln(args ...interface{}) {
	l.logln(SeverityWarning, args...)
}

// Info logs to the INFO severity logs.
func (l *Logger) Info(args ...interface{}) {
	l.log(SeverityInfo, args...)
}

// Infof logs to the INFO severity logs.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(SeverityInfo, format, args...)
}

// Infoln logs to the INFO severity logs.
func (l *Logger) Infoln(args ...interface{}) {
	l.logln(SeverityInfo, args...)
}

// Debug logs to the DEBUG severity logs.
func (l *Logger) Debug(args ...interface{}) {
	l.log(SeverityDebug, args...)
}

// Debugf logs to the DEBUG severity logs.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(SeverityDebug, format, args...)
}

// Debugln logs to the DEBUG severity logs.
func (l *Logger) Debugln(args ...interface{}) {
	l.logln(SeverityDebug, args...)
}

// Print logs to logs has Logger's print severity.
func (l *Logger) Print(args ...interface{}) {
	l.log(l.printSeverity, args...)
}

// Printf logs to logs has Logger's print severity.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.logf(l.printSeverity, format, args...)
}

// Println logs to logs has Logger's print severity.
func (l *Logger) Println(args ...interface{}) {
	l.logln(l.printSeverity, args...)
}

// SetOutput sets the Logger's output.
func (l *Logger) SetOutput(out Output) {
	l.mu.Lock()
	l.out = out
	l.mu.Unlock()
}

// SetSeverity sets the Logger's severity. If severity is invalid, it sets SeverityInfo.
func (l *Logger) SetSeverity(severity Severity) {
	l.mu.Lock()
	if !severity.IsValid() {
		severity = SeverityInfo
	}
	l.severity = severity
	l.mu.Unlock()
}

// SetVerbose sets the Logger's verbose.
func (l *Logger) SetVerbose(verbose Verbose) {
	l.mu.Lock()
	l.verbose = verbose
	l.mu.Unlock()
}

// V clones the Logger with given verbosity.
func (l *Logger) V(verbosity Verbose) *Logger {
	ln := l.clone()
	ln.verbosity = verbosity
	return ln
}

// SetPrintSeverity sets the Logger's severity level which is using with Print methods.
// If printSeverity is invalid, it sets SeverityInfo. By default, SeverityInfo.
func (l *Logger) SetPrintSeverity(printSeverity Severity) {
	l.mu.Lock()
	if !printSeverity.IsValid() {
		printSeverity = SeverityInfo
	}
	l.printSeverity = printSeverity
	l.mu.Unlock()
}

// SetStackTraceSeverity sets the Logger's severity level which allows printing stack trace.
// If stackTraceSeverity is invalid, it sets SeverityNone. By default, SeverityNone.
func (l *Logger) SetStackTraceSeverity(stackTraceSeverity Severity) {
	l.mu.Lock()
	if !stackTraceSeverity.IsValid() {
		stackTraceSeverity = SeverityNone
	}
	l.stackTraceSeverity = stackTraceSeverity
	l.mu.Unlock()
}

// WithTime clones the Logger with given time.
func (l *Logger) WithTime(tm time.Time) *Logger {
	ln := l.clone()
	ln.tm = tm
	return ln
}

// WithFields clones the Logger with given fields.
func (l *Logger) WithFields(fields ...Field) *Logger {
	ln := l.clone()
	ln.fields = append(ln.fields, fields...)
	return ln
}

// WithFieldKeyVals clones the Logger with given key and values of Field.
func (l *Logger) WithFieldKeyVals(kvs ...interface{}) *Logger {
	n := len(kvs)/2
	fields := make(Fields, 0, n)
	for i := 0; i < n; i++ {
		j := i*2
		k, v := fmt.Sprintf("%v", kvs[j]) , kvs[j+1]
		fields = append(fields, Field{Key: k, Val: v})
	}
	return l.WithFields(fields...)
}

// WithFieldMap clones the Logger with given fieldMap.
func (l *Logger) WithFieldMap(fieldMap map[string]interface{}) *Logger {
	fields := make(Fields, 0, len(fieldMap))
	for k, v := range fieldMap {
		fields = append(fields, Field{Key: k, Val: v})
	}
	return l.WithFields(fields...)
}
