// Package xlog provides leveled and structured logging.
package xlog

import (
	"io"
	"os"
	"time"

	"github.com/goinsane/erf"
)

var (
	defaultLogger *Logger     = New(defaultOutput, SeverityInfo, 0)
	defaultOutput *TextOutput = NewTextOutput(os.Stderr)

	defaultPCSize = erf.DefaultPCSize
)

// DefaultLogger returns the default Logger.
func DefaultLogger() *Logger {
	return defaultLogger
}

// DefaultOutput returns the default Output as Output. Type of the default Output is TextOutput.
func DefaultOutput() Output {
	return defaultOutput
}

// Fatal logs to the FATAL severity logs to the default Logger, then calls os.Exit(1).
func Fatal(args ...interface{}) {
	defaultLogger.log(SeverityFatal, args...)
	os.Exit(1)
}

// Fatalf logs to the FATAL severity logs to the default Logger, then calls os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	defaultLogger.logf(SeverityFatal, format, args...)
	os.Exit(1)
}

// Fatalln logs to the FATAL severity logs to the default Logger, then calls os.Exit(1).
func Fatalln(args ...interface{}) {
	defaultLogger.logln(SeverityFatal, args...)
	os.Exit(1)
}

// Error logs to the ERROR severity logs to the default Logger.
func Error(args ...interface{}) {
	defaultLogger.log(SeverityError, args...)
}

// Errorf logs to the ERROR severity logs to the default Logger.
func Errorf(format string, args ...interface{}) {
	defaultLogger.logf(SeverityError, format, args...)
}

// Errorln logs to the ERROR severity logs to the default Logger.
func Errorln(args ...interface{}) {
	defaultLogger.logln(SeverityError, args...)
}

// Warning logs to the WARNING severity logs to the default Logger.
func Warning(args ...interface{}) {
	defaultLogger.log(SeverityWarning, args...)
}

// Warningf logs to the WARNING severity logs to the default Logger.
func Warningf(format string, args ...interface{}) {
	defaultLogger.logf(SeverityWarning, format, args...)
}

// Warningln logs to the WARNING severity logs to the default Logger.
func Warningln(args ...interface{}) {
	defaultLogger.logln(SeverityWarning, args...)
}

// Info logs to the INFO severity logs to the default Logger.
func Info(args ...interface{}) {
	defaultLogger.log(SeverityInfo, args...)
}

// Infof logs to the INFO severity logs to the default Logger.
func Infof(format string, args ...interface{}) {
	defaultLogger.logf(SeverityInfo, format, args...)
}

// Infoln logs to the INFO severity logs to the default Logger.
func Infoln(args ...interface{}) {
	defaultLogger.logln(SeverityInfo, args...)
}

// Debug logs to the DEBUG severity logs to the default Logger.
func Debug(args ...interface{}) {
	defaultLogger.log(SeverityDebug, args...)
}

// Debugf logs to the DEBUG severity logs to the default Logger.
func Debugf(format string, args ...interface{}) {
	defaultLogger.logf(SeverityDebug, format, args...)
}

// Debugln logs to the DEBUG severity logs to the default Logger.
func Debugln(args ...interface{}) {
	defaultLogger.logln(SeverityDebug, args...)
}

// Print logs to the default Logger.
func Print(args ...interface{}) {
	defaultLogger.log(defaultLogger.printSeverity, args...)
}

// Printf logs to the default Logger.
func Printf(format string, args ...interface{}) {
	defaultLogger.logf(defaultLogger.printSeverity, format, args...)
}

// Println logs to the default Logger.
func Println(args ...interface{}) {
	defaultLogger.logln(defaultLogger.printSeverity, args...)
}

// SetOutput sets the default Logger's output. By default, the default Output.
func SetOutput(output Output) {
	defaultLogger.SetOutput(output)
}

// SetSeverity sets the default Logger's severity. If severity is invalid, it sets SeverityInfo.
// By default, SeverityInfo.
func SetSeverity(severity Severity) {
	defaultLogger.SetSeverity(severity)
}

// SetVerbose sets the default Logger's verbose.
// By default, 0.
func SetVerbose(verbose Verbose) {
	defaultLogger.SetVerbose(verbose)
}

// SetFlags sets the default Logger's flags.
// By default, FlagDefault.
func SetFlags(flags Flag) {
	defaultLogger.SetFlags(flags)
}

// SetPrintSeverity sets the default Logger's severity level which is using with Print methods.
// If printSeverity is invalid, it sets SeverityInfo. By default, SeverityInfo.
func SetPrintSeverity(printSeverity Severity) {
	defaultLogger.SetPrintSeverity(printSeverity)
}

// SetStackTraceSeverity sets the default Logger's severity level which allows printing stack trace.
// If stackTraceSeverity is invalid, it sets SeverityNone. By default, SeverityNone.
func SetStackTraceSeverity(stackTraceSeverity Severity) {
	defaultLogger.SetStackTraceSeverity(stackTraceSeverity)
}

// V clones the default Logger if Logger's verbose is greater or equal than given verbosity. Otherwise returns nil.
func V(verbosity Verbose) *Logger {
	return defaultLogger.V(verbosity)
}

// WithPrefix clones the default Logger and adds given prefix to end of the underlying prefix.
func WithPrefix(args ...interface{}) *Logger {
	return defaultLogger.WithPrefix(args...)
}

// WithPrefixf clones the default Logger and adds given prefix to end of the underlying prefix.
func WithPrefixf(format string, args ...interface{}) *Logger {
	return defaultLogger.WithPrefixf(format, args...)
}

// WithTime clones the default Logger with given time.
func WithTime(tm time.Time) *Logger {
	return defaultLogger.WithTime(tm)
}

// WithFields clones the default Logger with given fields.
func WithFields(fields ...Field) *Logger {
	return defaultLogger.WithFields(fields...)
}

// WithFieldKeyVals clones the default Logger with given key and values of Field.
func WithFieldKeyVals(kvs ...interface{}) *Logger {
	return defaultLogger.WithFieldKeyVals(kvs...)
}

// SetOutputWriter sets the default Output's writer.
func SetOutputWriter(w io.Writer) {
	defaultOutput.SetWriter(w)
}

// SetOutputFlags sets the default Output's flags to override every single Log.Flags if the flags argument different than 0.
// By default, 0.
func SetOutputFlags(flags Flag) {
	defaultOutput.SetFlags(flags)
}

// Reset resets the default Logger and the default Output.
func Reset() {
	SetOutput(defaultOutput)
	SetSeverity(SeverityInfo)
	SetVerbose(0)
	SetFlags(FlagDefault)
	SetPrintSeverity(SeverityInfo)
	SetStackTraceSeverity(SeverityNone)
	SetOutputWriter(os.Stderr)
	SetOutputFlags(0)
}
