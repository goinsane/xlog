// Package xlog provides leveled and structured logging.
package xlog

import (
	"io"
	"os"
	"time"

	"github.com/goinsane/erf"
)

var (
	defaultLogger       = New(defaultOutput, SeverityInfo, 0)
	defaultOutput       = NewTextOutput(defaultOutputWriter)
	defaultOutputWriter = os.Stderr
	defaultPCSize       = erf.DefaultPCSize
)

// DefaultLogger returns the default Logger.
func DefaultLogger() *Logger {
	return defaultLogger
}

// DefaultOutput returns the default Output as TextOutput type.
func DefaultOutput() *TextOutput {
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

// Print logs a log which has the default Logger's print severity to the default Logger.
func Print(args ...interface{}) {
	defaultLogger.log(defaultLogger.printSeverity, args...)
}

// Printf logs a log which has the default Logger's print severity to the default Logger.
func Printf(format string, args ...interface{}) {
	defaultLogger.logf(defaultLogger.printSeverity, format, args...)
}

// Println logs a log which has the default Logger's print severity to the default Logger.
func Println(args ...interface{}) {
	defaultLogger.logln(defaultLogger.printSeverity, args...)
}

// SetOutput sets the default Logger's output.
// It returns the default Logger.
// By default, the default Output.
func SetOutput(output Output) *Logger {
	return defaultLogger.SetOutput(output)
}

// SetSeverity sets the default Logger's severity.
// If severity is invalid, it sets SeverityInfo.
// It returns the default Logger.
// By default, SeverityInfo.
func SetSeverity(severity Severity) *Logger {
	return defaultLogger.SetSeverity(severity)
}

// SetVerbose sets the default Logger's verbose.
// It returns the default Logger.
// By default, 0.
func SetVerbose(verbose Verbose) *Logger {
	return defaultLogger.SetVerbose(verbose)
}

// SetFlags sets flags of Log created by the default Logger.
// These flags don't affect the default Logger. The Logger set them directly into the Log.
// It returns the default Logger.
// By default, FlagDefault.
func SetFlags(flags Flag) *Logger {
	return defaultLogger.SetFlags(flags)
}

// SetPrintSeverity sets the default Logger's severity level which is using with Print methods.
// If printSeverity is invalid, it sets SeverityInfo.
// It returns the default Logger.
// By default, SeverityInfo.
func SetPrintSeverity(printSeverity Severity) *Logger {
	return defaultLogger.SetPrintSeverity(printSeverity)
}

// SetStackTraceSeverity sets the default Logger's severity level which saves stack trace into Log.
// If stackTraceSeverity is invalid, it sets SeverityNone.
// It returns the default Logger.
// By default, SeverityNone.
func SetStackTraceSeverity(stackTraceSeverity Severity) *Logger {
	return defaultLogger.SetStackTraceSeverity(stackTraceSeverity)
}

// V duplicates the default Logger if the default Logger's verbose is greater or equal to given verbosity, otherwise returns nil.
func V(verbosity Verbose) *Logger {
	return defaultLogger.V(verbosity)
}

// WithPrefix duplicates the default Logger and adds given prefix to end of the underlying prefix.
func WithPrefix(args ...interface{}) *Logger {
	return defaultLogger.WithPrefix(args...)
}

// WithPrefixf duplicates the default Logger and adds given prefix to end of the underlying prefix.
func WithPrefixf(format string, args ...interface{}) *Logger {
	return defaultLogger.WithPrefixf(format, args...)
}

// WithTime duplicates the default Logger with given time.
func WithTime(tm time.Time) *Logger {
	return defaultLogger.WithTime(tm)
}

// WithFields duplicates the default Logger with given fields.
func WithFields(fields ...Field) *Logger {
	return defaultLogger.WithFields(fields...)
}

// WithFieldKeyVals duplicates the default Logger with given key and values of Field.
func WithFieldKeyVals(kvs ...interface{}) *Logger {
	return defaultLogger.WithFieldKeyVals(kvs...)
}

// SetOutputWriter sets the default Output's writer.
// It returns the default Output as TextOutput type.
// By default, os.Stderr.
func SetOutputWriter(w io.Writer) *TextOutput {
	return defaultOutput.SetWriter(w)
}

// SetOutputFlags sets the default Output's flags to override every single Log.Flags if the argument flags different from 0.
// It returns the default Output as TextOutput type.
// By default, 0.
func SetOutputFlags(flags Flag) *TextOutput {
	return defaultOutput.SetFlags(flags)
}

// ErfError creates a new *erf.Erf by given arguments. It logs to the ERROR severity logs to the default Logger and returns the new *erf.Erf.
func ErfError(args ...interface{}) *erf.Erf {
	return defaultLogger.erfError(SeverityError, args...)
}

// ErfErrorf creates a new *erf.Erf by given arguments. It logs to the ERROR severity logs to the default Logger and the result to get the new *erf.Erf.
func ErfErrorf(format string, args ...interface{}) *loggerErfResult {
	return defaultLogger.erfErrorf(SeverityError, format, args...)
}

// ErfWarning creates a new *erf.Erf by given arguments. It logs to the WARNING severity logs to the default Logger and returns the new *erf.Erf.
func ErfWarning(args ...interface{}) *erf.Erf {
	return defaultLogger.erfError(SeverityWarning, args...)
}

// ErfWarningf creates a new *erf.Erf by given arguments. It logs to the WARNING severity logs to the default Logger and returns the result to get the new *erf.Erf.
func ErfWarningf(format string, args ...interface{}) *loggerErfResult {
	return defaultLogger.erfErrorf(SeverityWarning, format, args...)
}

// Reset resets the default Logger and the default Output.
func Reset() {
	SetOutput(defaultOutput)
	SetSeverity(SeverityInfo)
	SetVerbose(0)
	SetFlags(FlagDefault)
	SetPrintSeverity(SeverityInfo)
	SetStackTraceSeverity(SeverityNone)
	SetOutputWriter(defaultOutputWriter)
	SetOutputFlags(0)
}
