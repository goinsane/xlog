package xlog

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/goinsane/erf"
)

// Log carries the log.
type Log struct {
	Message     []byte
	Severity    Severity
	Verbosity   Verbose
	Time        time.Time
	Fields      Fields
	StackCaller erf.StackCaller
	StackTrace  *erf.StackTrace
	Flags       Flag
}

// Duplicate duplicates the Log.
func (l *Log) Duplicate() *Log {
	if l == nil {
		return nil
	}
	l2 := &Log{
		Message:     nil,
		Severity:    l.Severity,
		Verbosity:   l.Verbosity,
		Time:        l.Time,
		Fields:      l.Fields.Duplicate(),
		StackCaller: l.StackCaller,
		StackTrace:  l.StackTrace.Duplicate(),
		Flags:       l.Flags,
	}
	if l.Message != nil {
		l2.Message = make([]byte, len(l.Message))
		copy(l2.Message, l.Message)
	}
	return l2
}

// String is implementation of fmt.Stringer.
func (l *Log) String() string {
	return fmt.Sprintf("%s", l)
}

// Format is implementation of fmt.Formatter.
func (l *Log) Format(f fmt.State, verb rune) {
	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	switch verb {
	case 's', 'v':
		if l.Flags&(FlagDate|FlagTime|FlagMicroseconds) != 0 {
			tm := l.Time
			if l.Flags&FlagUTC != 0 {
				tm = tm.UTC()
			}
			b := make([]byte, 0, 128)
			if l.Flags&FlagDate != 0 {
				year, month, day := l.Time.Date()
				itoa(&b, year, 4)
				b = append(b, '/')
				itoa(&b, int(month), 2)
				b = append(b, '/')
				itoa(&b, day, 2)
				b = append(b, ' ')
			}
			if l.Flags&(FlagTime|FlagMicroseconds) != 0 {
				hour, min, sec := l.Time.Clock()
				itoa(&b, hour, 2)
				b = append(b, ':')
				itoa(&b, min, 2)
				b = append(b, ':')
				itoa(&b, sec, 2)
				if l.Flags&FlagMicroseconds != 0 {
					b = append(b, '.')
					itoa(&b, l.Time.Nanosecond()/1e3, 6)
				}
				b = append(b, ' ')
			}
			buf.Write(b)
		}

		if l.Flags&FlagSeverity != 0 {
			buf.WriteString(l.Severity.String())
			buf.WriteString(" - ")
		}

		var padding []byte
		if l.Flags&FlagPadding != 0 {
			padding = bytes.Repeat([]byte(" "), buf.Len())
		}

		if l.Flags&(FlagLongFunc|FlagShortFunc) != 0 {
			fn := "???"
			if l.StackCaller.Function != "" {
				fn = trimSrcPath(l.StackCaller.Function)
			}
			if l.Flags&FlagShortFunc != 0 {
				fn = trimDirs(fn)
			}
			buf.WriteString(fn)
			buf.WriteString("()")
			buf.WriteString(" - ")
		}

		if l.Flags&(FlagLongFile|FlagShortFile) != 0 {
			file, line := "???", 0
			if l.StackCaller.File != "" {
				file = trimSrcPath(l.StackCaller.File)
			}
			if l.StackCaller.Line > 0 {
				line = l.StackCaller.Line
			}
			if l.Flags&FlagShortFile != 0 {
				file = trimDirs(file)
			}
			buf.WriteString(file)
			buf.WriteRune(':')
			b := make([]byte, 0, 128)
			itoa(&b, line, -1)
			buf.Write(b)
			buf.WriteString(" - ")
		}

		for idx, line := range bytes.Split(l.Message, []byte("\n")) {
			if idx > 0 {
				buf.Write(padding)
			}
			buf.Write(line)
			buf.WriteRune('\n')
		}

		if l.Flags&FlagFields != 0 && len(l.Fields) > 0 {
			fields := l.Fields.Duplicate()
			sort.Sort(fields)
			buf.WriteRune('\t')
			for idx, field := range fields {
				if idx > 0 {
					buf.WriteRune(' ')
				}
				buf.WriteString(fmt.Sprintf("%s=%q", field.Key, fmt.Sprintf("%v", field.Value)))
			}
			buf.WriteRune('\n')
		}

		if l.Flags&FlagStackTrace != 0 && l.StackTrace != nil {
			buf.WriteString(fmt.Sprintf("%+1.1v", l.StackTrace))
			buf.WriteRune('\n')
		}
	default:
		return
	}
	_, _ = f.Write(buf.Bytes())
}

// MarshalText is implementation of encoding.TextMarshaler.
func (l *Log) MarshalText() (text []byte, err error) {
	f := &fmtState{}
	l.Format(f, 's')
	return f.Buffer, nil
}
