//go:build examples
// +build examples

package main

import (
	"os"
	"time"

	"github.com/goinsane/xlog"
)

var (
	testTime, _ = time.ParseInLocation("2006-01-02T15:04:05", "2010-11-12T13:14:15", time.Local)
)

func main() {
	// reset xlog for previous changes if it is running in go test.
	xlog.Reset()
	// change writer of default output to stdout from stderr.
	xlog.SetOutputWriter(os.Stdout)

	// log by Severity.
	// default severity is SeverityInfo.
	// default verbose is 0.
	xlog.Debug("this is debug log. but it won't be shown.")
	xlog.Info("this is info log.")
	xlog.Warning("this is warning log.")
	xlog.V(1).Error("this is error log, verbosity 1. but it won't be shown.")

	// SetSeverity()
	xlog.SetSeverity(xlog.SeverityDebug)
	xlog.Debug("this is debug log. it will now be shown.")

	// SetVerbose() and V()
	xlog.SetVerbose(1)
	xlog.V(1).Error("this is error log, verbosity 1. it will now be shown.")
	xlog.V(2).Warning("this is warning log, verbosity 2. it won't be shown.")

	// SetFlags()
	// default flags is FlagDefault.
	xlog.SetFlags(xlog.FlagDefault | xlog.FlagShortFile)
	xlog.Info("this is info log. you can see file name and line in this log.")

	// log using Print.
	// default print severity is SeverityInfo.
	xlog.Print("this log will be shown as info log.")

	// SetPrintSeverity()
	xlog.SetPrintSeverity(xlog.SeverityWarning)
	xlog.Print("this log will now be shown as warning log.")

	// SetStackTraceSeverity()
	// default stack trace severity is none.
	xlog.SetStackTraceSeverity(xlog.SeverityWarning)
	xlog.Warning("this is warning log. you can see stack trace end of this log.")
	xlog.Error("this is error log. you can still see stack trace end of this log.")
	xlog.Info("this is info log. stack trace won't be shown end of this log.")

	// WithPrefix()
	xlog.WithPrefix("prefix1").Warning("this is warning log with prefix 'prefix1'.")
	xlog.WithPrefix("prefix1").WithPrefix("prefix2").Error("this is error log with both of prefixes 'prefix1' and 'prefix2'.")

	// WithTime()
	xlog.WithTime(testTime).Info("this is info log with custom time.")

	// WithFieldKeyVals()
	xlog.WithFieldKeyVals("key1", "val1", "key2", "val2", "key3", "val3", "key1", "val1-2", "key2", "val2-2").Info("this is info log with several fields.")
}
