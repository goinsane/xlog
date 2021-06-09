// Package xlog provides leveled and structured logging.
package xlog

import (
	"io"
	"os"
	"time"
)

var (
	defLogger *Logger     = New(defOutput, SeverityInfo, 0)
	defOutput *TextOutput = NewTextOutput(os.Stderr, OutputFlagDefault)
)

// DefaultLogger returns the default logger.
func DefaultLogger() *Logger {
	return defLogger
}

// DefaultOutput returns the default output as Output. Type of default output is TextOutput.
func DefaultOutput() Output {
	return defOutput
}

// Fatal logs to the FATAL severity logs to the default logger, then calls os.Exit(1).
func Fatal(args ...interface{}) {
	defLogger.log(SeverityFatal, args...)
	os.Exit(1)
}

// Fatalf logs to the FATAL severity logs to the default logger, then calls os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	defLogger.logf(SeverityFatal, format, args...)
	os.Exit(1)
}

// Fatalln logs to the FATAL severity logs to the default logger, then calls os.Exit(1).
func Fatalln(args ...interface{}) {
	defLogger.logln(SeverityFatal, args...)
	os.Exit(1)
}

// Error logs to the ERROR severity logs to the default logger.
func Error(args ...interface{}) {
	defLogger.log(SeverityError, args...)
}

// Errorf logs to the ERROR severity logs to the default logger.
func Errorf(format string, args ...interface{}) {
	defLogger.logf(SeverityError, format, args...)
}

// Errorln logs to the ERROR severity logs to the default logger.
func Errorln(args ...interface{}) {
	defLogger.logln(SeverityError, args...)
}

// Warning logs to the WARNING severity logs to the default logger.
func Warning(args ...interface{}) {
	defLogger.log(SeverityWarning, args...)
}

// Warningf logs to the WARNING severity logs to the default logger.
func Warningf(format string, args ...interface{}) {
	defLogger.logf(SeverityWarning, format, args...)
}

// Warningln logs to the WARNING severity logs to the default logger.
func Warningln(args ...interface{}) {
	defLogger.logln(SeverityWarning, args...)
}

// Info logs to the INFO severity logs to the default logger.
func Info(args ...interface{}) {
	defLogger.log(SeverityInfo, args...)
}

// Infof logs to the INFO severity logs to the default logger.
func Infof(format string, args ...interface{}) {
	defLogger.logf(SeverityInfo, format, args...)
}

// Infoln logs to the INFO severity logs to the default logger.
func Infoln(args ...interface{}) {
	defLogger.logln(SeverityInfo, args...)
}

// Debug logs to the DEBUG severity logs to the default logger.
func Debug(args ...interface{}) {
	defLogger.log(SeverityDebug, args...)
}

// Debugf logs to the DEBUG severity logs to the default logger.
func Debugf(format string, args ...interface{}) {
	defLogger.logf(SeverityDebug, format, args...)
}

// Debugln logs to the DEBUG severity logs to the default logger.
func Debugln(args ...interface{}) {
	defLogger.logln(SeverityDebug, args...)
}

// Print logs to the default logger.
func Print(args ...interface{}) {
	defLogger.log(defLogger.printSeverity, args...)
}

// Printf logs to the default logger.
func Printf(format string, args ...interface{}) {
	defLogger.logf(defLogger.printSeverity, format, args...)
}

// Println logs to the default logger.
func Println(args ...interface{}) {
	defLogger.logln(defLogger.printSeverity, args...)
}

// SetOutput sets the default logger's output. By default, the default output.
func SetOutput(out Output) {
	defLogger.SetOutput(out)
}

// SetSeverity sets the default logger's severity. If severity is invalid, it sets SeverityInfo.
// By default, SeverityInfo.
func SetSeverity(severity Severity) {
	defLogger.SetSeverity(severity)
}

// SetVerbose sets the default logger's verbose. By default, 0.
func SetVerbose(verbose Verbose) {
	defLogger.SetVerbose(verbose)
}

// SetPrintSeverity sets the default logger's severity level which is using with Print functions.
// If printSeverity is invalid, it sets SeverityInfo. By default, SeverityInfo.
func SetPrintSeverity(printSeverity Severity) {
	defLogger.SetPrintSeverity(printSeverity)
}

// SetStackTraceSeverity sets the default logger's severity level which allows printing stack trace.
// If stackTraceSeverity is invalid, it sets SeverityNone. By default, SeverityNone.
func SetStackTraceSeverity(stackTraceSeverity Severity) {
	defLogger.SetStackTraceSeverity(stackTraceSeverity)
}

// V clones the default logger if Logger's verbose is greater or equal than given verbosity. Otherwise returns nil.
func V(verbosity Verbose) *Logger {
	return defLogger.V(verbosity)
}

// WithPrefix clones the default Logger and adds given prefix to end of the underlying prefix.
func WithPrefix(args ...interface{}) *Logger {
	return defLogger.WithPrefix(args...)
}

// WithPrefixf clones the default Logger and adds given prefix to end of the underlying prefix.
func WithPrefixf(format string, args ...interface{}) *Logger {
	return defLogger.WithPrefixf(format, args...)
}

// WithTime clones the default logger with given time.
func WithTime(tm time.Time) *Logger {
	return defLogger.WithTime(tm)
}

// WithFields clones the default logger with given fields.
func WithFields(fields ...Field) *Logger {
	return defLogger.WithFields(fields...)
}

// WithFieldKeyVals clones default logger with given key and values of Field.
func WithFieldKeyVals(kvs ...interface{}) *Logger {
	return defLogger.WithFieldKeyVals(kvs...)
}

// SetOutputWriter sets the default output writer.
func SetOutputWriter(w io.Writer) {
	defOutput.SetWriter(w)
}

// SetOutputFlags sets the default output flags.
func SetOutputFlags(flags OutputFlag) {
	defOutput.SetFlags(flags)
}

// SetOutputPadding sets custom padding of the default output. If padding is empty-string, padding is filled by first line of log.
func SetOutputPadding(padding string) {
	defOutput.SetPadding(padding)
}

// Reset resets default logger and output options.
func Reset() {
	SetOutput(defOutput)
	SetSeverity(SeverityInfo)
	SetVerbose(0)
	SetPrintSeverity(SeverityInfo)
	SetStackTraceSeverity(SeverityNone)
	SetOutputWriter(os.Stderr)
	SetOutputFlags(OutputFlagDefault)
}
