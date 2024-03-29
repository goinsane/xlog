package xlog

import (
	"fmt"
	"strings"
)

// Field is type of field.
type Field struct {
	Key   string
	Value interface{}

	mark interface{}
}

func (f *Field) GetMark() interface{} {
	return f.mark
}

// Fields is slice of fields.
type Fields []Field

// Duplicate duplicates the Fields.
func (f Fields) Duplicate() Fields {
	if f == nil {
		return nil
	}
	f2 := make(Fields, 0, len(f))
	for i := range f {
		f2 = append(f2, f[i])
	}
	return f2
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

type FieldMarkErf struct {
	No    int
	Index int
}

func (m *FieldMarkErf) String() string {
	return fmt.Sprintf("%d:%d", m.No, m.Index)
}
