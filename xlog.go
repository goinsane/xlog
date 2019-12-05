// Package xlog provides leveled and structured logging.
package xlog

import (
	"bytes"
	"fmt"
	"go/build"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// Verbose is type of verbose level.
type Verbose uint16

// Field is type of field.
type Field struct {
	Key string
	Val interface{}
}

// Fields is type of fields.
type Fields []Field

// Clone clones Fields.
func (f Fields) Clone() Fields {
	if f == nil {
		return nil
	}
	result := make(Fields, 0, len(f))
	for i := range f {
		result = append(result, f[i])
	}
	return result
}

// Len is implementation of sort.Interface
func (f Fields) Len() int {
	return len(f)
}

// Less is implementation of sort.Interface
func (f Fields) Less(i, j int) bool {
	return strings.Compare(f[i].Key, f[j].Key) < 0
}

// Swap is implementation of sort.Interface
func (f Fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// Callers is a type of stack callers.
type Callers []uintptr

var (
	defLogger *Logger     = New(defOutput, SeverityInfo, 0)
	defOutput *TextOutput = NewTextOutput(os.Stdout, OutputFlagDefault)
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

func trimSrcpath(s string) string {
	var r string
	r = strings.TrimPrefix(s, build.Default.GOROOT+"/src/")
	if r != s {
		return r
	}
	r = strings.TrimPrefix(s, build.Default.GOPATH+"/src/")
	if r != s {
		return r
	}
	return s
}

// CallersToStackTrace generates stack trace output from stack callers.
func CallersToStackTrace(callers Callers, padding []byte) []byte {
	if callers == nil {
		return nil
	}
	frames := runtime.CallersFrames(callers)
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	for {
		frame, more := frames.Next()
		buf.Write(padding)
		buf.WriteString(fmt.Sprintf("%s\n", frame.Function))
		buf.Write(padding)
		buf.WriteString(fmt.Sprintf("\t%s:%d\n", trimSrcpath(frame.File), frame.Line))
		if !more {
			break
		}
	}
	return buf.Bytes()
}

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

// V clones the default logger with given verbosity.
func V(verbosity Verbose) *Logger {
	return defLogger.V(verbosity)
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

// Reset resets default logger and output options.
func Reset() {
	SetOutput(defOutput)
	SetSeverity(SeverityInfo)
	SetVerbose(0)
	SetPrintSeverity(SeverityInfo)
	SetStackTraceSeverity(SeverityNone)
	SetOutputWriter(os.Stdout)
	SetOutputFlags(OutputFlagDefault)
}
