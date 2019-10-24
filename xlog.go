// Package xlog provides leveled and structured logging.
package xlog

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

// Verbose is type of verbose level.
type Verbose uint16

// Fields is type of fields.
type Fields map[string]interface{}

// Callers is a type of stack callers.
type Callers []uintptr

var (
	defLogger *Logger = New(defOutput, SeverityInfo, 0)
	defOutput *TextOutput  = NewTextOutput(os.Stdout, OutputFlagDefault)
)

func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// CallersToStackTrace generates stack trace output from stack callers.
func CallersToStackTrace(callers Callers, padding []byte) []byte {
	frames := runtime.CallersFrames(callers)
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	for {
		frame, more := frames.Next()
		buf.Write(padding)
		buf.WriteString(fmt.Sprintf("%s\n", frame.Function))
		buf.Write(padding)
		buf.WriteString(fmt.Sprintf("\t%s:%d\n", frame.File, frame.Line))
		if !more {
			break
		}
	}
	return buf.Bytes()
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

// V clones the default logger with given verbosity.
func V(verbosity Verbose) *Logger {
	return defLogger.V(verbosity)
}

// WithTime clones the default logger with given time.
func WithTime(tm time.Time) *Logger {
	return defLogger.WithTime(tm)
}

// WithFields clones the default logger with given fields.
func WithFields(fields Fields) *Logger {
	return defLogger.WithFields(fields)
}

// SetOutputFlags sets the default output flags.
func SetOutputFlags(flags OutputFlag) {
	defOutput.SetFlags(flags)
}

// SetOutputStackTraceSeverity sets the default output severity level which allows printing stack trace.
func SetOutputStackTraceSeverity(stackTraceSeverity Severity) {
	defOutput.SetStackTraceSeverity(stackTraceSeverity)
}
