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
	printSeverity      Severity
	stackTraceSeverity Severity
	prefix             string
	verbosity          Verbose
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
		printSeverity:      SeverityInfo,
		stackTraceSeverity: SeverityNone,
	}
}

// Duplicate duplicates the Logger.
func (l *Logger) Duplicate() *Logger {
	l.mu.RLock()
	ln := &Logger{
		out:                l.out,
		severity:           l.severity,
		verbose:            l.verbose,
		printSeverity:      l.printSeverity,
		stackTraceSeverity: l.stackTraceSeverity,
		prefix:             l.prefix,
		verbosity:          l.verbosity,
		tm:                 l.tm,
		fields:             l.fields.Duplicate(),
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
		msg := &Message{}
		messageLen := len(l.prefix) + len(message)
		msg.Msg = make([]byte, 0, messageLen+1)
		msg.Msg = append(msg.Msg, l.prefix...)
		msg.Msg = append(msg.Msg, message...)
		if messageLen == 0 || msg.Msg[messageLen-1] != '\n' {
			msg.Msg = append(msg.Msg, '\n')
		}
		msg.Severity = severity
		msg.Verbosity = l.verbosity
		msg.Tm = l.tm
		if msg.Tm.IsZero() {
			msg.Tm = time.Now()
		}
		msg.Caller, _, _, _ = runtime.Caller(3)
		if f := runtime.FuncForPC(msg.Caller); f != nil {
			msg.Func = trimSrcPath(f.Name())
			msg.File, msg.Line = f.FileLine(msg.Caller)
			msg.File = trimSrcPath(msg.File)
		}
		msg.Fields = l.fields.Duplicate()
		if l.stackTraceSeverity >= severity {
			msg.Callers = make(Callers, 32)
			msg.Callers = msg.Callers[:runtime.Callers(4, msg.Callers)]
		}
		l.out.Log(msg)
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
	if l == nil {
		return
	}
	l.log(l.printSeverity, args...)
}

// Printf logs to logs has Logger's print severity.
func (l *Logger) Printf(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.logf(l.printSeverity, format, args...)
}

// Println logs to logs has Logger's print severity.
func (l *Logger) Println(args ...interface{}) {
	if l == nil {
		return
	}
	l.logln(l.printSeverity, args...)
}

// SetOutput sets the Logger's output.
func (l *Logger) SetOutput(out Output) {
	if l == nil {
		return
	}
	l.mu.Lock()
	l.out = out
	l.mu.Unlock()
}

// SetSeverity sets the Logger's severity. If severity is invalid, it sets SeverityInfo.
func (l *Logger) SetSeverity(severity Severity) {
	if l == nil {
		return
	}
	l.mu.Lock()
	if !severity.IsValid() {
		severity = SeverityInfo
	}
	l.severity = severity
	l.mu.Unlock()
}

// SetVerbose sets the Logger's verbose.
func (l *Logger) SetVerbose(verbose Verbose) {
	if l == nil {
		return
	}
	l.mu.Lock()
	l.verbose = verbose
	l.mu.Unlock()
}

// SetPrintSeverity sets the Logger's severity level which is using with Print methods.
// If printSeverity is invalid, it sets SeverityInfo. By default, SeverityInfo.
func (l *Logger) SetPrintSeverity(printSeverity Severity) {
	if l == nil {
		return
	}
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
	if l == nil {
		return
	}
	l.mu.Lock()
	if !stackTraceSeverity.IsValid() {
		stackTraceSeverity = SeverityNone
	}
	l.stackTraceSeverity = stackTraceSeverity
	l.mu.Unlock()
}

// V clones the Logger if Logger's verbose is greater or equal than given verbosity. Otherwise returns nil.
func (l *Logger) V(verbosity Verbose) *Logger {
	if l == nil {
		return nil
	}
	if !(l.verbose >= verbosity) {
		return nil
	}
	ln := l.Duplicate()
	ln.verbosity = verbosity
	return ln
}

// WithPrefix clones the Logger and adds given prefix to end of the underlying prefix.
func (l *Logger) WithPrefix(args ...interface{}) *Logger {
	if l == nil {
		return nil
	}
	ln := l.Duplicate()
	ln.prefix += fmt.Sprint(args...) + ": "
	return ln
}

// WithPrefixf clones the Logger and adds given prefix to end of the underlying prefix.
func (l *Logger) WithPrefixf(format string, args ...interface{}) *Logger {
	if l == nil {
		return nil
	}
	ln := l.Duplicate()
	ln.prefix += fmt.Sprintf(format, args...) + ": "
	return ln
}

// WithTime clones the Logger with given time.
func (l *Logger) WithTime(tm time.Time) *Logger {
	if l == nil {
		return nil
	}
	ln := l.Duplicate()
	ln.tm = tm
	return ln
}

// WithFields clones the Logger with given fields.
func (l *Logger) WithFields(fields ...Field) *Logger {
	if l == nil {
		return nil
	}
	ln := l.Duplicate()
	ln.fields = append(ln.fields, fields...)
	return ln
}

// WithFieldKeyVals clones the Logger with given key and values of Field.
func (l *Logger) WithFieldKeyVals(kvs ...interface{}) *Logger {
	if l == nil {
		return nil
	}
	n := len(kvs) / 2
	fields := make(Fields, 0, n)
	for i := 0; i < n; i++ {
		j := i * 2
		k, v := fmt.Sprintf("%v", kvs[j]), kvs[j+1]
		fields = append(fields, Field{Key: k, Val: v})
	}
	return l.WithFields(fields...)
}

// WithFieldMap clones the Logger with given fieldMap.
func (l *Logger) WithFieldMap(fieldMap map[string]interface{}) *Logger {
	if l == nil {
		return nil
	}
	fields := make(Fields, 0, len(fieldMap))
	for k, v := range fieldMap {
		fields = append(fields, Field{Key: k, Val: v})
	}
	return l.WithFields(fields...)
}
