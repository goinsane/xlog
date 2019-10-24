package xlog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Output interface {
	Log([]byte, Severity, Verbose, time.Time, Fields, Callers)
}

type OutputFlag int

const (
	OutputFlagDate         = OutputFlag(1 << iota)                                // the date in the local time zone: 2009/01/23
	OutputFlagTime                                                                // the time in the local time zone: 01:23:23
	OutputFlagMicroseconds                                                        // microsecond resolution: 01:23:23.123123
	OutputFlagUTC                                                                 // use UTC rather than the local time zone
	OutputFlagSeverity                                                            // severity level
	OutputFlagPadding                                                             // use padding multiple lines
	OutputFlagLongFile                                                            // full file name and line number: /a/b/c/d.go:23
	OutputFlagShortFile                                                           // final file name element and line number: d.go:23
	OutputFlagStackTrace                                                          // print stack trace
	OutputFlagDefault      = OutputFlagDate | OutputFlagTime | OutputFlagSeverity // initial values for the default logger

)

type TextOutput struct {
	mu                 sync.Mutex
	w                  io.Writer
	bw                 *bufio.Writer
	flags              OutputFlag
	stackTraceSeverity Severity
}

func NewTextOutput(w io.Writer, flags OutputFlag) *TextOutput {
	return &TextOutput{
		w:                  w,
		bw:                 bufio.NewWriter(w),
		flags:              flags,
		stackTraceSeverity: SeverityInfo,
	}
}

func (o *TextOutput) Log(msg []byte, severity Severity, verbose Verbose, tm time.Time, fields Fields, callers Callers) {
	var err error
	o.mu.Lock()
	defer o.mu.Unlock()

	buf := make([]byte, 128)
	padLen := 0

	buf = buf[:0]
	if o.flags&(OutputFlagDate|OutputFlagTime|OutputFlagMicroseconds) != 0 {
		if o.flags&OutputFlagUTC != 0 {
			tm = tm.UTC()
		}
		if o.flags&OutputFlagDate != 0 {
			year, month, day := tm.Date()
			itoa(&buf, year, 4)
			buf = append(buf, '/')
			itoa(&buf, int(month), 2)
			buf = append(buf, '/')
			itoa(&buf, day, 2)
			buf = append(buf, ' ')
		}
		if o.flags&(OutputFlagTime|OutputFlagMicroseconds) != 0 {
			hour, min, sec := tm.Clock()
			itoa(&buf, hour, 2)
			buf = append(buf, ':')
			itoa(&buf, min, 2)
			buf = append(buf, ':')
			itoa(&buf, sec, 2)
			if o.flags&OutputFlagMicroseconds != 0 {
				buf = append(buf, '.')
				itoa(&buf, tm.Nanosecond()/1e3, 6)
			}
			buf = append(buf, ' ')
		}
	}
	if o.flags&OutputFlagSeverity != 0 {
		buf = append(buf, fmt.Sprintf("%7s: ", severity.String())...)
	}
	if o.flags&OutputFlagPadding != 0 {
		padLen = len(buf)
	}
	_, err = o.bw.Write(buf)
	if err != nil {
		return
	}

	padding := strings.Repeat(" ", padLen)

	for i := 0; len(msg) > 0; i++ {
		if i > 0 {
			_, err = o.bw.WriteString(padding)
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
		_, err = o.bw.Write(msg[:idx])
		if err != nil {
			return
		}
		msg = msg[idx:]
	}

	if o.flags&(OutputFlagLongFile|OutputFlagShortFile) != 0 {
		buf = buf[:0]
		buf = append(buf, "\tFile: "...)
		file, line := "???", 0
		if len(callers) > 0 {
			f := runtime.FuncForPC(callers[0])
			file, line = f.FileLine(callers[0])
		}
		if o.flags&OutputFlagShortFile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		buf = append(buf, file...)
		buf = append(buf, ':')
		itoa(&buf, line, -1)
		buf = append(buf, '\n')
		_, err = o.bw.Write(buf)
		if err != nil {
			return
		}
	}

	if o.flags&OutputFlagStackTrace != 0 && severity <= o.stackTraceSeverity {
		buf = buf[:0]
		//buf = append(buf, "\tStack Trace: \n"...)
		buf = append(buf, CallersToStackTrace(callers, []byte("\t"))...)
		//buf = append(buf, '\n')
		_, err = o.bw.Write(buf)
		if err != nil {
			return
		}
	}

	o.bw.Flush()
}

func (o *TextOutput) SetFlags(flags OutputFlag) {
	o.mu.Lock()
	o.flags = flags
	o.mu.Unlock()
}

func (o *TextOutput) SetStackTraceSeverity(stackTraceSeverity Severity) {
	o.mu.Lock()
	if !stackTraceSeverity.IsValid() {
		stackTraceSeverity = SeverityInfo
	}
	o.stackTraceSeverity = stackTraceSeverity
	o.mu.Unlock()
}
