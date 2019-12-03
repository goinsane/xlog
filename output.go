package xlog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// Output is an interface for Logger output. All of Output implementations must be
// concurrent-safe.
type Output interface {
	Log([]byte, Severity, Verbose, time.Time, Fields, Callers)
}

// OutputFlag is type of output flag.
type OutputFlag int

const (
	// OutputFlagDate prints the date in the local time zone: 2009/01/23
	OutputFlagDate = OutputFlag(1 << iota)

	// OutputFlagTime prints the time in the local time zone: 01:23:23
	OutputFlagTime

	// OutputFlagMicroseconds prints microsecond resolution: 01:23:23.123123
	OutputFlagMicroseconds

	// OutputFlagUTC uses UTC rather than the local time zone
	OutputFlagUTC

	// OutputFlagSeverity prints severity level
	OutputFlagSeverity

	// OutputFlagPadding prints padding with multiple lines
	OutputFlagPadding

	// OutputFlagFields prints fields
	OutputFlagFields

	// OutputFlagLongFile prints full file name and line number: /a/b/c/d.go:23
	OutputFlagLongFile

	// OutputFlagShortFile prints final file name element and line number: d.go:23
	OutputFlagShortFile

	// OutputFlagStackTrace prints stack trace
	OutputFlagStackTrace

	// OutputFlagDefault holds initial values for the default logger
	OutputFlagDefault = OutputFlagDate | OutputFlagTime | OutputFlagSeverity | OutputFlagFields | OutputFlagStackTrace
)

// TextOutput is an implementation of Output by writing texts to io.Writer w.
type TextOutput struct {
	mu    sync.Mutex
	w     io.Writer
	bw    *bufio.Writer
	flags OutputFlag
}

// NewTextOutput creates a new TextOutput.
func NewTextOutput(w io.Writer, flags OutputFlag) *TextOutput {
	return &TextOutput{
		w:     w,
		bw:    bufio.NewWriter(w),
		flags: flags,
	}
}

// Log implementes Output.Log
func (o *TextOutput) Log(msg []byte, severity Severity, verbose Verbose, tm time.Time, fields Fields, callers Callers) {
	var err error
	o.mu.Lock()
	defer o.mu.Unlock()

	defer o.bw.Flush()

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
		buf = append(buf, severity.String()...)
		buf = append(buf, ": "...)
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
			idx = len(msg)
		} else {
			idx++
		}
		_, err = o.bw.Write(msg[:idx])
		if err != nil {
			return
		}
		msg = msg[idx:]
	}

	if len(fields) > 0 && o.flags&OutputFlagFields != 0 {
		fields2 := fields.Clone()
		sort.Sort(fields2)
		buf = buf[:0]
		buf = append(buf, "\tFields: "...)
		for _, f := range fields2 {
			buf = append(buf, fmt.Sprintf("%s=%q ", f.Key, fmt.Sprintf("%v", f.Val))...)
		}
		buf = append(buf[:len(buf)-1], '\n')
		_, err = o.bw.Write(buf)
		if err != nil {
			return
		}
	}

	if len(callers) > 0 {
		if o.flags&(OutputFlagLongFile|OutputFlagShortFile) != 0 {
			buf = buf[:0]
			buf = append(buf, "\tFile: "...)
			file, line := "???", 0
			f := runtime.FuncForPC(callers[0])
			if f != nil {
				file, line = f.FileLine(callers[0])
				file = trimSrcpath(file)
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
		if o.flags&OutputFlagStackTrace != 0 {
			buf = buf[:0]
			buf = append(buf, CallersToStackTrace(callers, []byte("\t"))...)
			_, err = o.bw.Write(buf)
			if err != nil {
				return
			}
		}
	}
}

// SetWriter sets output writer.
func (o *TextOutput) SetWriter(w io.Writer) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.bw.Flush()
	o.w = w
	o.bw = bufio.NewWriter(w)
}

// SetFlags sets output flags.
func (o *TextOutput) SetFlags(flags OutputFlag) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.flags = flags
}
