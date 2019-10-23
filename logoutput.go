package xlog

import (
	"runtime"
	"time"
)

type LogOutput interface {
	Log([]byte, time.Time, Fields, Severity, *runtime.Frames)
	Flush()
}

var defLogOutput LogOutput
