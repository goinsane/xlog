package xlog_test

import (
	"os"

	"github.com/goinsane/xlog"
)

func ExampleSimple() {
	xlog.SetVerbose(2)

	xlog.Info("this is info log, verbosity 0")
	xlog.V(0).Warning("this is warning log, verbosity 0")
	xlog.V(1).Warning("this is warning log, verbosity 1")
	xlog.V(2).Error("this is error log, verbosity 2")
	xlog.V(3).Error("this is error log, verbosity 3. it will not be shown")
	xlog.Debug("this is debug log, verbosity 0")
}

func ExampleDetail() {
	output := xlog.NewTextOutput(os.Stdout, xlog.OutputFlagDefault | xlog.OutputFlagLongFile | xlog.OutputFlagPadding)
	logger := xlog.New(output, xlog.SeverityInfo, 2)
	logger.SetStackTraceSeverity(xlog.SeverityError)

	logger.Info("this is info log, \nverbosity 0")
	logger.V(0).Warning("this is warning log, verbosity 0")
	logger.V(1).Warning("this is warning log, verbosity 1")
	logger.V(2).Error("this is error log, verbosity 2")
	logger.V(3).Error("this is error log, verbosity 3. it will not be shown")
	logger.Debug("this is debug log, verbosity 0. it will not be shown")
}

func ExampleTest() {
	logger := xlog.New(xlog.NewTextOutput(os.Stdout, xlog.OutputFlagSeverity), xlog.SeverityInfo, 2)

	logger.Info("this is info log, verbosity 0")
	logger.V(0).Warning("this is warning log, verbosity 0")
	logger.V(1).Warning("this is warning log, verbosity 1")
	logger.V(2).Error("this is error log, verbosity 2")
	logger.V(3).Error("this is error log, verbosity 3. it will not be shown")
	logger.Debug("this is debug log, verbosity 0. it will not be shown")

	logger.SetSeverity(xlog.SeverityDebug)
	logger.Debug("this is debug log, verbosity 0. it will be shown, now")

	// Output:
	// INFO: this is info log, verbosity 0
	// WARNING: this is warning log, verbosity 0
	// WARNING: this is warning log, verbosity 1
	// ERROR: this is error log, verbosity 2
	// DEBUG: this is debug log, verbosity 0. it will be shown, now
}
