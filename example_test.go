package xlog_test

import (
	"os"

	"github.com/goinsane/xlog"
)

func Example() {
	xlog.SetSeverity(xlog.SeverityDebug)
	xlog.SetVerbose(2)

	xlog.Info("this is info log, verbosity 0")
	xlog.V(0).Warning("this is warning log, verbosity 0")
	xlog.V(1).Warning("this is warning log, verbosity 1")
	xlog.V(2).Error("this is error log, verbosity 2")
	xlog.V(3).Error("this is error log, verbosity 3. it will not be shown")
	xlog.Debug("this is debug log, verbosity 0")

	// Output:
}

func ExampleLogger() {
	logger := xlog.New(xlog.NewTextOutput(os.Stdout, xlog.OutputFlagSeverity), xlog.SeverityInfo, 2)

	logger.Info("this is info log, verbosity 0")
	logger.V(0).Warning("this is warning log, verbosity 0")
	logger.V(1).Warning("this is warning log, verbosity 1")
	logger.V(2).Error("this is error log, verbosity 2")
	logger.V(3).Error("this is error log, verbosity 3. it will not be shown")
	logger.Debug("this is debug log, verbosity 0. it will not be shown")

	// Output:
	// INFO: this is info log, verbosity 0
	// WARNING: this is warning log, verbosity 0
	// WARNING: this is warning log, verbosity 1
	// ERROR: this is error log, verbosity 2
}

func ExampleTextOutput() {
	output := xlog.NewTextOutput(os.Stdout, xlog.OutputFlagSeverity|xlog.OutputFlagStackTrace)
	logger := xlog.New(output, xlog.SeverityInfo, 2)

	output.SetStackTraceSeverity(xlog.SeverityWarning)

	logger.Info("this is info log, verbosity 0")
	logger.V(0).Warning("this is warning log, verbosity 0")
	logger.V(1).Warning("this is warning log, verbosity 1")
	logger.V(2).Error("this is error log, verbosity 2")
	logger.V(3).Error("this is error log, verbosity 3. it will not be shown")
	logger.Debug("this is debug log, verbosity 0. it will not be shown")

	// Output:
}
