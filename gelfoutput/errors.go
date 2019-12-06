package gelfoutput

import (
	"errors"
)

var (
	// ErrUnknownGelfWriterType defines 'unknown gelf writer type' error
	ErrUnknownGelfWriterType = errors.New("unknown gelf writer type")
)
