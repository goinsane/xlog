package xlog

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/goinsane/erf"
)

// Logger provides a logger for leveled and structured logging.
type Logger struct {
	mu                 sync.RWMutex
	output             Output
	severity           Severity
	verbose            Verbose
	printSeverity      Severity
	stackTraceSeverity Severity
	prefix             string
	verbosity          Verbose
	time               time.Time
	fields             Fields
}

// New creates a new Logger. If severity is invalid, it sets SeverityInfo.
func New(output Output, severity Severity, verbose Verbose) *Logger {
	if !severity.IsValid() {
		severity = SeverityInfo
	}
	return &Logger{
		output:             output,
		severity:           severity,
		verbose:            verbose,
		printSeverity:      SeverityInfo,
		stackTraceSeverity: SeverityNone,
	}
}

// Duplicate duplicates the Logger.
func (l *Logger) Duplicate() *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l2 := &Logger{
		output:             l.output,
		severity:           l.severity,
		verbose:            l.verbose,
		printSeverity:      l.printSeverity,
		stackTraceSeverity: l.stackTraceSeverity,
		prefix:             l.prefix,
		verbosity:          l.verbosity,
		time:               l.time,
		fields:             l.fields.Duplicate(),
	}
	return l2
}

func (l *Logger) out(severity Severity, message string) {
	if l == nil {
		return
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.output != nil && l.severity >= severity && l.verbose >= l.verbosity {
		messageLen := len(l.prefix) + len(message)
		log := &Log{
			Message:   make([]byte, 0, messageLen+1),
			Severity:  severity,
			Verbosity: l.verbosity,
			Time:      l.time,
			Fields:    l.fields.Duplicate(),
		}
		log.Message = append(log.Message, l.prefix...)
		log.Message = append(log.Message, message...)
		if messageLen == 0 || log.Message[messageLen-1] != '\n' {
			log.Message = append(log.Message, '\n')
		}
		if log.Time.IsZero() {
			log.Time = time.Now()
		}
		log.StackCaller = erf.NewStackTrace(erf.PC(1, 5)...).Caller(0)
		if l.stackTraceSeverity >= severity {
			log.StackTrace = erf.NewStackTrace(erf.PC(erf.DefaultPCSize, 5)...)
		}
		l.output.Log(log)
	}
}

func (l *Logger) log(severity Severity, args ...interface{}) {
	l.out(severity, fmt.Sprint(args...))
}

func (l *Logger) logf(severity Severity, format string, args ...interface{}) {
	l.out(severity, fmt.Sprintf(format, args...))
}

func (l *Logger) logln(severity Severity, args ...interface{}) {
	l.out(severity, fmt.Sprintln(args...))
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
func (l *Logger) SetOutput(output Output) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = output
}

// SetSeverity sets the Logger's severity. If severity is invalid, it sets SeverityInfo.
func (l *Logger) SetSeverity(severity Severity) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if !severity.IsValid() {
		severity = SeverityInfo
	}
	l.severity = severity
}

// SetVerbose sets the Logger's verbose.
func (l *Logger) SetVerbose(verbose Verbose) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.verbose = verbose
}

// SetPrintSeverity sets the Logger's severity level which is using with Print methods.
// If printSeverity is invalid, it sets SeverityInfo. By default, SeverityInfo.
func (l *Logger) SetPrintSeverity(printSeverity Severity) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if !printSeverity.IsValid() {
		printSeverity = SeverityInfo
	}
	l.printSeverity = printSeverity
}

// SetStackTraceSeverity sets the Logger's severity level which allows printing stack trace.
// If stackTraceSeverity is invalid, it sets SeverityNone. By default, SeverityNone.
func (l *Logger) SetStackTraceSeverity(stackTraceSeverity Severity) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if !stackTraceSeverity.IsValid() {
		stackTraceSeverity = SeverityNone
	}
	l.stackTraceSeverity = stackTraceSeverity
}

// V duplicates the Logger if Logger's verbose is greater or equal than given verbosity. Otherwise returns nil.
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

// WithPrefix duplicates the Logger and adds given prefix to end of the underlying prefix.
func (l *Logger) WithPrefix(args ...interface{}) *Logger {
	if l == nil {
		return nil
	}
	ln := l.Duplicate()
	ln.prefix += fmt.Sprint(args...) + ": "
	return ln
}

// WithPrefixf duplicates the Logger and adds given prefix to end of the underlying prefix.
func (l *Logger) WithPrefixf(format string, args ...interface{}) *Logger {
	if l == nil {
		return nil
	}
	ln := l.Duplicate()
	ln.prefix += fmt.Sprintf(format, args...) + ": "
	return ln
}

// WithTime duplicates the Logger with given time.
func (l *Logger) WithTime(tm time.Time) *Logger {
	if l == nil {
		return nil
	}
	ln := l.Duplicate()
	ln.time = tm
	return ln
}

// WithFields duplicates the Logger with given fields.
func (l *Logger) WithFields(fields ...Field) *Logger {
	if l == nil {
		return nil
	}
	ln := l.Duplicate()
	ln.fields = append(ln.fields, fields...)
	return ln
}

// WithFieldKeyVals duplicates the Logger with given key and values of Field.
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

// WithFieldMap duplicates the Logger with given fieldMap.
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
