package xlog_test

import (
	"time"

	"github.com/goinsane/xlog"
)

func Example() {
	xlog.Reset()

	// set output flags for testable results
	xlog.SetOutputFlags(xlog.OutputFlagSeverity | xlog.OutputFlagPadding)

	xlog.Debug("this is debug log, verbosity 0. it will not be shown.")
	xlog.Info("this is info log, verbosity 0.")
	xlog.V(1).Info("this is info log, verbosity 1. it will not be shown.")
	xlog.Print("this is info log, verbosity 0 caused by Print().")

	xlog.SetSeverity(xlog.SeverityDebug)
	xlog.Debug("this is debug log, verbosity 0.")

	xlog.SetVerbose(1)
	xlog.Info("this is info log, verbosity 0.")
	xlog.V(1).Info("this is info log, verbosity 1.")

	xlog.SetPrintSeverity(xlog.SeverityWarning)
	xlog.Print("this is warning log, verbosity 0 caused by Print().")

	xlog.Error("this is error log,\nverbosity 0.")

	xlog.SetOutputFlags(xlog.OutputFlagSeverity | xlog.OutputFlagPadding | xlog.OutputFlagFields)
	xlog.WithFieldKeyVals("key2", "val2", "key3", "val3", "key2", "val2-2", "key1", "val1").Warning("this is warning log,\nverbosity 0 with fields.")

	xlog.SetOutputFlags(xlog.OutputFlagDefault | xlog.OutputFlagUTC)
	tm, _ := time.Parse(time.RFC3339, "2019-11-13T21:56:24+00:00")
	xlog.WithTime(tm).Info("this is info log, verbosity 0.")

	// Output:
	// INFO: this is info log, verbosity 0.
	// INFO: this is info log, verbosity 0 caused by Print().
	// DEBUG: this is debug log, verbosity 0.
	// INFO: this is info log, verbosity 0.
	// INFO: this is info log, verbosity 1.
	// WARNING: this is warning log, verbosity 0 caused by Print().
	// ERROR: this is error log,
	//        verbosity 0.
	// WARNING: this is warning log,
	//          verbosity 0 with fields.
	// 	Fields: key1="val1" key2="val2" key2="val2-2" key3="val3"
	// 2019/11/13 21:56:24 INFO: this is info log, verbosity 0.
}
