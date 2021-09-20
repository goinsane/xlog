// +build examples

package main

import (
	"os"
	"time"

	"github.com/goinsane/xlog"
)

func main() {
	// reset xlog for previous changes
	xlog.Reset()
	xlog.SetFlags(xlog.FlagSeverity)
	xlog.SetOutputWriter(os.Stdout)

	xlog.Debug("this is debug log, verbosity 0. it will not be shown.")
	xlog.Info("this is info log, verbosity 0.")
	xlog.Warning("this is warning log, verbosity 0.")
	xlog.Error("this is error log, verbosity 0.")
	xlog.Print("this is info log, verbosity 0 caused by Print().")
	xlog.V(1).Info("this is info log, verbosity 1. it will not be shown.")

	xlog.SetSeverity(xlog.SeverityDebug)
	xlog.Debug("this is debug log, verbosity 0.")

	xlog.SetVerbose(1)
	xlog.V(0).Info("this is info log, verbosity 0.")
	xlog.V(1).Info("this is info log, verbosity 1.")
	xlog.V(2).Info("this is info log, verbosity 2. it will not be shown.")

	xlog.SetPrintSeverity(xlog.SeverityWarning)
	xlog.Print("this is warning log, verbosity 0 caused by Print().")

	xlog.Warning("this is warning log, verbosity 0.\nwithout padding.")
	xlog.SetFlags(xlog.FlagSeverity | xlog.FlagPadding)
	xlog.Warning("this is warning log, verbosity 0.\nwith padding.")

	xlog.SetFlags(xlog.FlagSeverity | xlog.FlagPadding | xlog.FlagFields)
	xlog.WithFieldKeyVals("key1", "val1", "key2", "val2", "key3", "val3", "key1", "val1-2", "key2", "val2-2").Error("this is error log, verbosity 0.\nwith padding.\nwith fields.")

	xlog.SetFlags(xlog.FlagDefault)
	tm, _ := time.ParseInLocation("2006-01-02T15:04:05", "2019-11-13T21:56:24", time.Local)
	xlog.WithTime(tm).Info("this is info log, verbosity 0.")

	// Output:
	// INFO - this is info log, verbosity 0.
	// WARNING - this is warning log, verbosity 0.
	// ERROR - this is error log, verbosity 0.
	// INFO - this is info log, verbosity 0 caused by Print().
	// DEBUG - this is debug log, verbosity 0.
	// INFO - this is info log, verbosity 0.
	// INFO - this is info log, verbosity 1.
	// WARNING - this is warning log, verbosity 0 caused by Print().
	// WARNING - this is warning log, verbosity 0.
	// without padding.
	// WARNING - this is warning log, verbosity 0.
	//           with padding.
	// ERROR - this is error log, verbosity 0.
	//         with padding.
	//         with fields.
	// 	key1="val1" key2="val2" key3="val3" key1="val1-2" key2="val2-2"
	// 2019/11/13 21:56:24 INFO - this is info log, verbosity 0.
}
