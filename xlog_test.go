package xlog_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/goinsane/xlog"
)

func BenchmarkLogger_Info(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}

func BenchmarkLogger_Infof(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("%s", "benchmark")
	}
}

func BenchmarkLogger_Infoln(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infoln("benchmark")
	}
}

func BenchmarkLogger_Info_withStackTrace(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	logger.SetStackTraceSeverity(xlog.SeverityInfo)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}

func BenchmarkLogger_Info_withFlagLongFunc(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	logger.SetFlags(xlog.FlagDefault | xlog.FlagLongFunc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}

func BenchmarkLogger_Info_withFlagShortFunc(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	logger.SetFlags(xlog.FlagDefault | xlog.FlagShortFunc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}

func BenchmarkLogger_Info_withFlagLongFile(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	logger.SetFlags(xlog.FlagDefault | xlog.FlagLongFile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}

func BenchmarkLogger_Info_withFlagShortFile(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	logger.SetFlags(xlog.FlagDefault | xlog.FlagShortFile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}

func BenchmarkLogger_V(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.V(1)
	}
}

func BenchmarkLogger_WithPrefix(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithPrefix("prefix")
	}
}

func BenchmarkLogger_WithPrefixf(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithPrefixf("%s", "prefix")
	}
}

func BenchmarkLogger_WithTime(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithTime(testTime)
	}
}

func BenchmarkLogger_WithFieldKeyVals(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFieldKeyVals("key1", "value1")
	}
}

func ExampleSetSeverity() {
	reset()
	xlog.SetSeverity(xlog.SeverityDebug)
	xlog.Debug("this is debug log, verbosity 0.")
	xlog.Info("this is info log, verbosity 0.")
	xlog.Warning("this is warning log, verbosity 0.")

	// Output:
	// DEBUG - this is debug log, verbosity 0.
	// INFO - this is info log, verbosity 0.
	// WARNING - this is warning log, verbosity 0.
}

func ExampleSetVerbose() {
	reset()
	xlog.SetVerbose(2)
	xlog.V(0).Debug("this is debug log, verbosity 0. it won't be shown.")
	xlog.V(1).Info("this is info log, verbosity 1.")
	xlog.V(2).Warning("this is warning log, verbosity 2.")
	xlog.V(3).Error("this is error log, verbosity 3. it won't be shown.")

	// Output:
	// INFO - this is info log, verbosity 1.
	// WARNING - this is warning log, verbosity 2.
}

func ExampleSetFlags() {
	reset()
	xlog.SetFlags(0)
	xlog.Info("this is info log, verbosity 0.")

	// Output:
	// this is info log, verbosity 0.
}

func ExampleWithTime() {
	reset()
	xlog.SetFlags(xlog.FlagDefault)
	xlog.WithTime(testTime).Info("this is info log, verbosity 0.")

	// Output:
	// 2010/11/12 13:14:15 INFO - this is info log, verbosity 0.
}

func ExampleLogger() {
	logger := xlog.New(xlog.NewTextOutput(os.Stdout), xlog.SeverityInfo, 2)
	logger.SetFlags(xlog.FlagSeverity)

	logger.Info("this is info log, verbosity 0.")
	logger.V(0).Info("this is info log, verbosity 0.")
	logger.V(1).Warning("this is warning log, verbosity 1.")
	logger.V(2).Error("this is error log, verbosity 2.")
	logger.V(3).Error("this is error log, verbosity 3. it won't be shown.")
	logger.Debug("this is debug log, verbosity 0. it won't be shown.")

	// Output:
	// INFO - this is info log, verbosity 0.
	// INFO - this is info log, verbosity 0.
	// WARNING - this is warning log, verbosity 1.
	// ERROR - this is error log, verbosity 2.
}

func Example_example1() {
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

var (
	testTime, _ = time.ParseInLocation("2006-01-02T15:04:05", "2010-11-12T13:14:15", time.Local)
)

// reset resets xlog for previous changes
func reset() {
	xlog.Reset()
	xlog.SetFlags(xlog.FlagDefault & ^xlog.FlagDate & ^xlog.FlagTime)
	xlog.SetOutputWriter(os.Stdout)
}
