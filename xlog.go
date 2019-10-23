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
	defLogger    *Logger   = New(defLogOutput, SeverityInfo, 0)
	defLogOutput LogOutput = NewTextLogOutput(os.Stdout, LogOutputFlagDefault)
)

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
