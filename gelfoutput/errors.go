package gelfoutput

import (
	"errors"
)

var (
	ErrUnknownGelfWriterType = errors.New("unknown gelf writer type")
)
