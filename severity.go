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
	sSeverityFatal   = "fatal"
	sSeverityError   = "error"
	sSeverityWarning = "warning"
	sSeverityInfo    = "info"
	sSeverityDebug   = "debug"
	sSeverityUnknown = "unknown"
)

var sSeverities = []string{sSeverityFatal, sSeverityError, sSeverityWarning, sSeverityInfo, sSeverityDebug}

func ParseSeverity(s string) Severity {
	s = strings.ToLower(s)
	if s == sSeverityFatal {
		return SeverityFatal
	}
	if s == sSeverityError || s == "err" {
		return SeverityError
	}
	if s == sSeverityWarning || s == "warn" {
		return SeverityWarning
	}
	if s == sSeverityInfo {
		return SeverityInfo
	}
	if s == sSeverityDebug || s == "dbg" {
		return SeverityDebug
	}
	return SeverityUnknown
}
