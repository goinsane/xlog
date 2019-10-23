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
	Flush()
}

type TextLogOutput struct {
	mu      sync.RWMutex
	w       io.Writer
	bw      *bufio.Writer
	padding string
}

func NewTextLogOutput(w io.Writer) *TextLogOutput {
	return &TextLogOutput{
		w:  w,
		bw: bufio.NewWriter(w),
	}
}

func (lo *TextLogOutput) Log(msg []byte, severity Severity, verbose Verbose, tm time.Time, fields Fields, callers *runtime.Frames) {
	var err error
	lo.mu.RLock()
	padding := []byte(lo.padding)
	lo.mu.RUnlock()
	for i := 0; len(msg) > 0; i++ {
		if i > 0 && len(padding) > 0 {
			_, err = lo.bw.Write(padding)
			if err != nil {
				break
			}
		}
		idx := bytes.IndexByte(msg, '\n')
		if idx < 0 {
			idx = len(msg)
		} else {
			idx++
		}
		_, err = lo.bw.Write(msg[:idx])
		if err != nil {
			break
		}
		msg = msg[idx:]
	}
	lo.bw.Flush()
}

func (lo *TextLogOutput) Flush() {
	lo.bw.Flush()
}

func (lo *TextLogOutput) SetPadding(padding string) {
	lo.mu.Lock()
	lo.padding = padding
	lo.mu.Unlock()
}
