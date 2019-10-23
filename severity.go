package xlog

import "strings"

type Severity int

func (sv Severity) String() string {
	idx := int(sv)
	if idx < len(sSeverities) {
		return sSeverities[idx]
	}
	return sSeverityUnknown
}

func (sv Severity) IsValid() bool {
	k := int(sv)
	return k >= 0 && k < len(severities)
}

const (
	SeverityFatal = Severity(iota)
	SeverityError
	SeverityWarning
	SeverityInfo
	SeverityDebug
)

var severities = []Severity{SeverityFatal, SeverityError, SeverityWarning, SeverityInfo, SeverityDebug}

const SeverityUnknown = -1

const (
	sSeverityFatal   = "FATAL"
	sSeverityError   = "ERROR"
	sSeverityWarning = "WARNING"
	sSeverityInfo    = "INFO"
	sSeverityDebug   = "DEBUG"
	sSeverityUnknown = "UNKNOWN"
)

var sSeverities = []string{sSeverityFatal, sSeverityError, sSeverityWarning, sSeverityInfo, sSeverityDebug}

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
	return SeverityUnknown
}
