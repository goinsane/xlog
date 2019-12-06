package xlog

import "strings"

// Severity is type of severity level.
type Severity int

// String returns severity name by string.
func (sv Severity) String() string {
	if sv.IsValid() {
		return sSeverities[int(sv)]
	}
	return sSeverityNone
}

// IsValid checks Severity value is valid.
func (sv Severity) IsValid() bool {
	k := int(sv)
	return k >= 0 && k < len(severities)
}

const (
	// SeverityFatal is fatal severity level
	SeverityFatal = Severity(iota)

	// SeverityError is error severity level
	SeverityError

	// SeverityWarning is warning severity level
	SeverityWarning

	// SeverityInfo is info severity level
	SeverityInfo

	// SeverityDebug is debug severity level
	SeverityDebug
)

var severities = []Severity{SeverityFatal, SeverityError, SeverityWarning, SeverityInfo, SeverityDebug}

// SeverityNone is none or unknown severity level
const SeverityNone = -1

const (
	sSeverityFatal   = "FATAL"
	sSeverityError   = "ERROR"
	sSeverityWarning = "WARNING"
	sSeverityInfo    = "INFO"
	sSeverityDebug   = "DEBUG"
	sSeverityNone    = "NONE"
)

var sSeverities = []string{sSeverityFatal, sSeverityError, sSeverityWarning, sSeverityInfo, sSeverityDebug}

// ParseSeverity parses severity name. If it fails, returns SeverityNone.
func ParseSeverity(s string) Severity {
	s = strings.ToUpper(s)
	if s == sSeverityFatal {
		return SeverityFatal
	}
	if s == sSeverityError || s == "ERR" {
		return SeverityError
	}
	if s == sSeverityWarning || s == "WARN" {
		return SeverityWarning
	}
	if s == sSeverityInfo {
		return SeverityInfo
	}
	if s == sSeverityDebug || s == "DBG" {
		return SeverityDebug
	}
	return SeverityNone
}
