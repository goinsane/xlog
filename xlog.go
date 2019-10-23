package xlog

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

type Verbose uint32

type Fields map[string]interface{}

var (
	defLogger    *Logger   = New(defLogOutput, SeverityInfo, 0)
	defLogOutput LogOutput = NewTextLogOutput(os.Stdout)
)

func FramesToStackTrace(frames *runtime.Frames, padding []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 256))
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
