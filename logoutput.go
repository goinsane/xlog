package xlog

import (
	"bufio"
	"bytes"
	"io"
	"runtime"
	"sync"
	"time"
)

type LogOutput interface {
	Log([]byte, Severity, Verbose, time.Time, Fields, *runtime.Frames)
}

type TextLogOutput struct {
	mu               sync.RWMutex
	logMu            sync.Mutex
	w                io.Writer
	bw               *bufio.Writer
	padding          []byte
	reportStackTrace bool
}

func NewTextLogOutput(w io.Writer) *TextLogOutput {
	return &TextLogOutput{
		w:  w,
		bw: bufio.NewWriter(w),
	}
}

func (lo *TextLogOutput) Log(msg []byte, severity Severity, verbose Verbose, tm time.Time, fields Fields, frames *runtime.Frames) {
	var err error
	lo.logMu.Lock()
	defer lo.logMu.Unlock()
	lo.mu.RLock()
	padding := lo.padding
	reportStackTrace := lo.reportStackTrace
	lo.mu.RUnlock()
	for i := 0; len(msg) > 0; i++ {
		if i > 0 && len(padding) > 0 {
			_, err = lo.bw.Write(padding)
			if err != nil {
				return
			}
		}
		idx := bytes.IndexByte(msg, '\n')
		if idx < 0 {
			msg = append(msg, '\n')
			idx = len(msg) + 1
		} else {
			idx++
		}
		_, err = lo.bw.Write(msg[:idx])
		if err != nil {
			return
		}
		msg = msg[idx:]
	}
	if reportStackTrace {
		lo.bw.Write(FramesToStackTrace(frames, padding))
	}
	lo.bw.Flush()
}

func (lo *TextLogOutput) SetPadding(padding string) {
	lo.mu.Lock()
	lo.padding = []byte(padding)
	lo.mu.Unlock()
}

func (lo *TextLogOutput) SetReportStackTrace(reportStackTrace bool) {
	lo.mu.Lock()
	lo.reportStackTrace = reportStackTrace
	lo.mu.Unlock()
}
