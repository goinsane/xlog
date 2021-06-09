package xlog

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// Verbose is type of verbose level.
type Verbose int

// Field is type of field.
type Field struct {
	Key string
	Val interface{}
}

// Fields is type of fields.
type Fields []Field

// Duplicate duplicates the Fields.
func (f Fields) Duplicate() Fields {
	if f == nil {
		return nil
	}
	result := make(Fields, 0, len(f))
	for i := range f {
		result = append(result, f[i])
	}
	return result
}

// Len is implementation of sort.Interface
func (f Fields) Len() int {
	return len(f)
}

// Less is implementation of sort.Interface
func (f Fields) Less(i, j int) bool {
	return strings.Compare(f[i].Key, f[j].Key) < 0
}

// Swap is implementation of sort.Interface
func (f Fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// Message carries log message.
type Message struct {
	Msg       []byte
	Severity  Severity
	Verbosity Verbose
	Tm        time.Time
	Caller    uintptr
	Func      string
	File      string
	Line      int
	Fields    Fields
	Callers   Callers
}

// Duplicate duplicates the Message.
func (msg *Message) Duplicate() *Message {
	if msg == nil {
		return nil
	}
	result := &Message{}
	*result = *msg
	result.Msg = make([]byte, len(msg.Msg))
	copy(result.Msg, msg.Msg)
	result.Fields = msg.Fields.Duplicate()
	result.Callers = msg.Callers.Clone()
	return result
}

// Callers is a type of stack callers.
type Callers []uintptr

func (c Callers) Clone() Callers {
	if c == nil {
		return nil
	}
	result := make(Callers, 0, len(c))
	for i := range c {
		result = append(result, c[i])
	}
	return result
}

// ToStackTrace generates stack trace output from stack callers.
func (c Callers) ToStackTrace(padding []byte) []byte {
	frames := runtime.CallersFrames(c)
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	for {
		frame, more := frames.Next()
		buf.Write(padding)
		buf.WriteString(fmt.Sprintf("%s()\n", trimSrcPath(frame.Function)))
		buf.Write(padding)
		buf.WriteString(fmt.Sprintf("\t%s:%d\n", trimSrcPath(frame.File), frame.Line))
		if !more {
			break
		}
	}
	return buf.Bytes()
}
