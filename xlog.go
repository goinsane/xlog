package xlog

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

type Verbose uint16

type Fields map[string]interface{}

type Callers []uintptr

var (
	defLogger    *Logger   = New(defLogOutput, SeverityInfo, 0)
	defLogOutput LogOutput = NewTextLogOutput(os.Stdout, LogOutputFlagDefault)
)

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
