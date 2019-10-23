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

type LogOutput interface {
	Log([]byte, Severity, Verbose, time.Time, Fields, Callers)
}

type LogOutputFlag int

const (
	LogOutputFlagDate         = LogOutputFlag(1 << iota)                                      // the date in the local time zone: 2009/01/23
	LogOutputFlagTime                                                                         // the time in the local time zone: 01:23:23
	LogOutputFlagMicroseconds                                                                 // microsecond resolution: 01:23:23.123123
	LogOutputFlagUTC                                                                          // use UTC rather than the local time zone
	LogOutputFlagSeverity                                                                     // severity level
	LogOutputFlagPadding                                                                      // use padding multiple lines
	LogOutputFlagLongFile                                                                     // full file name and line number: /a/b/c/d.go:23
	LogOutputFlagShortFile                                                                    // final file name element and line number: d.go:23
	LogOutputFlagStackTrace                                                                   // print stack trace
	LogOutputFlagDefault      = LogOutputFlagDate | LogOutputFlagTime | LogOutputFlagSeverity // initial values for the default logger

)

type TextLogOutput struct {
	mu                 sync.Mutex
	w                  io.Writer
	bw                 *bufio.Writer
	flags              LogOutputFlag
	stackTraceSeverity Severity
}

func NewTextLogOutput(w io.Writer, flags LogOutputFlag) *TextLogOutput {
	return &TextLogOutput{
		w:                  w,
		bw:                 bufio.NewWriter(w),
		flags:              flags,
		stackTraceSeverity: SeverityInfo,
	}
}

func (lo *TextLogOutput) Log(msg []byte, severity Severity, verbose Verbose, tm time.Time, fields Fields, callers Callers) {
	var err error
	lo.mu.Lock()
	defer lo.mu.Unlock()

	buf := make([]byte, 128)
	padLen := 0

	buf = buf[:0]
	if lo.flags&(LogOutputFlagDate|LogOutputFlagTime|LogOutputFlagMicroseconds) != 0 {
		if lo.flags&LogOutputFlagUTC != 0 {
			tm = tm.UTC()
		}
		if lo.flags&LogOutputFlagDate != 0 {
			year, month, day := tm.Date()
			itoa(&buf, year, 4)
			buf = append(buf, '/')
			itoa(&buf, int(month), 2)
			buf = append(buf, '/')
			itoa(&buf, day, 2)
			buf = append(buf, ' ')
		}
		if lo.flags&(LogOutputFlagTime|LogOutputFlagMicroseconds) != 0 {
			hour, min, sec := tm.Clock()
			itoa(&buf, hour, 2)
			buf = append(buf, ':')
			itoa(&buf, min, 2)
			buf = append(buf, ':')
			itoa(&buf, sec, 2)
			if lo.flags&LogOutputFlagMicroseconds != 0 {
				buf = append(buf, '.')
				itoa(&buf, tm.Nanosecond()/1e3, 6)
			}
			buf = append(buf, ' ')
		}
	}
	if lo.flags&LogOutputFlagSeverity != 0 {
		buf = append(buf, fmt.Sprintf("%7s: ", severity.String())...)
	}
	if lo.flags&LogOutputFlagPadding != 0 {
		padLen = len(buf)
	}
	_, err = lo.bw.Write(buf)
	if err != nil {
		return
	}

	padding := strings.Repeat(" ", padLen)

	for i := 0; len(msg) > 0; i++ {
		if i > 0 {
			_, err = lo.bw.WriteString(padding)
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

	if lo.flags&(LogOutputFlagLongFile|LogOutputFlagShortFile) != 0 {
		buf = buf[:0]
		buf = append(buf, "\tFile: "...)
		file, line := "???", 0
		if len(callers) > 0 {
			f := runtime.FuncForPC(callers[0])
			file, line = f.FileLine(callers[0])
		}
		if lo.flags&LogOutputFlagShortFile != 0 {
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
		_, err = lo.bw.Write(buf)
		if err != nil {
			return
		}
	}

	if lo.flags&LogOutputFlagStackTrace != 0 && severity <= lo.stackTraceSeverity {
		buf = buf[:0]
		//buf = append(buf, "\tStack Trace: \n"...)
		buf = append(buf, CallersToStackTrace(callers, []byte("\t"))...)
		//buf = append(buf, '\n')
		_, err = lo.bw.Write(buf)
		if err != nil {
			return
		}
	}

	lo.bw.Flush()
}

func (lo *TextLogOutput) SetFlags(flags LogOutputFlag) {
	lo.mu.Lock()
	lo.flags = flags
	lo.mu.Unlock()
}

func (lo *TextLogOutput) SetStackTraceSeverity(stackTraceSeverity Severity) {
	lo.mu.Lock()
	if !stackTraceSeverity.IsValid() {
		stackTraceSeverity = SeverityInfo
	}
	lo.stackTraceSeverity = stackTraceSeverity
	lo.mu.Unlock()
}
