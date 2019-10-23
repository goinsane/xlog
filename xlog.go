package xlog

import "os"

type Verbose uint32

type Fields map[string]interface{}

var (
	defLogger    *Logger   = New(defLogOutput, SeverityInfo, 0)
	defLogOutput LogOutput = NewTextLogOutput(os.Stdout)
)
