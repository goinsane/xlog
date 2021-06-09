package xlog

import (
	"strings"
	"time"

	"github.com/goinsane/erf"
)

// Verbose is type of verbose level.
type Verbose int

// Field is type of field.
type Field struct {
	Key string
	Val interface{}
}

// Fields is slice of fields.
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

// Len is implementation of sort.Interface.
func (f Fields) Len() int {
	return len(f)
}

// Less is implementation of sort.Interface.
func (f Fields) Less(i, j int) bool {
	return strings.Compare(f[i].Key, f[j].Key) < 0
}

// Swap is implementation of sort.Interface.
func (f Fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// Message carries log message.
type Message struct {
	Msg         []byte
	Severity    Severity
	Verbosity   Verbose
	Time        time.Time
	Fields      Fields
	StackCaller erf.StackCaller
	StackTrace  *erf.StackTrace
}

// Duplicate duplicates the Message.
func (m *Message) Duplicate() *Message {
	if m == nil {
		return nil
	}
	m2 := new(Message)
	*m2 = *m
	m2.Msg = make([]byte, len(m.Msg))
	copy(m2.Msg, m.Msg)
	m2.Fields = m.Fields.Duplicate()
	m2.StackTrace = m.StackTrace.Duplicate()
	return m2
}
