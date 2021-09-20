package xlog

import (
	"strings"

	"github.com/goinsane/erf"
)

var (
	ErrUnknownSeverity = erf.New("unknown severity")
)

// Severity describes severity level of Log.
type Severity int

const (
	// SeverityNone is none or unspecified severity level
	SeverityNone Severity = iota

	// SeverityFatal is fatal severity level
	SeverityFatal

	// SeverityError is error severity level
	SeverityError

	// SeverityWarning is warning severity level
	SeverityWarning

	// SeverityInfo is info severity level
	SeverityInfo

	// SeverityDebug is debug severity level
	SeverityDebug
)

// String is implementation of fmt.Stringer.
func (s Severity) String() string {
	text, _ := s.MarshalText()
	return string(text)
}

// IsValid returns whether value is valid.
func (s Severity) IsValid() bool {
	return s.CheckValid() == nil
}

// CheckValid returns error for invalid value.
func (s Severity) CheckValid() error {
	_, err := s.MarshalText()
	return err
}

// MarshalText is implementation of encoding.TextMarshaler.
func (s Severity) MarshalText() (text []byte, err error) {
	var str string
	switch s {
	case SeverityNone:
		str = "NONE"
	case SeverityFatal:
		str = "FATAL"
	case SeverityError:
		str = "ERROR"
	case SeverityWarning:
		str = "WARNING"
	case SeverityInfo:
		str = "INFO"
	case SeverityDebug:
		str = "DEBUG"
	default:
		return nil, erf.Wrap(ErrUnknownSeverity)
	}
	return []byte(str), nil
}

// UnmarshalText is implementation of encoding.UnmarshalText.
func (s *Severity) UnmarshalText(text []byte) error {
	switch str := strings.ToUpper(string(text)); str {
	case "NONE", "NON", "NA":
		*s = SeverityNone
	case "FATAL", "FTL":
		*s = SeverityFatal
	case "ERROR", "ERR":
		*s = SeverityError
	case "WARNING", "WRN", "WARN":
		*s = SeverityWarning
	case "INFO", "INF":
		*s = SeverityInfo
	case "DEBUG", "DBG":
		*s = SeverityDebug
	default:
		return erf.Wrap(ErrUnknownSeverity)
	}
	return nil
}
