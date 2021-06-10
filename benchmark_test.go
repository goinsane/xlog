package xlog_test

import (
	"io/ioutil"
	"testing"

	"github.com/goinsane/xlog"
)

func BenchmarkLogger(b *testing.B) {
	output := xlog.NewTextOutput(ioutil.Discard, xlog.FlagDefault)
	logger := xlog.New(output, xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}

func BenchmarkLoggerStackTrace(b *testing.B) {
	output := xlog.NewTextOutput(ioutil.Discard, xlog.FlagDefault)
	logger := xlog.New(output, xlog.SeverityInfo, 0)
	logger.SetStackTraceSeverity(xlog.SeverityInfo)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}

func BenchmarkLoggerV(b *testing.B) {
	output := xlog.NewTextOutput(ioutil.Discard, xlog.FlagDefault)
	logger := xlog.New(output, xlog.SeverityInfo, 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.V(1).Info("benchmark")
	}
}

func BenchmarkOutputFlagShortFunc(b *testing.B) {
	output := xlog.NewTextOutput(ioutil.Discard, xlog.FlagDefault| xlog.FlagShortFunc)
	logger := xlog.New(output, xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}

func BenchmarkOutputFlagShortFile(b *testing.B) {
	output := xlog.NewTextOutput(ioutil.Discard, xlog.FlagDefault| xlog.FlagShortFile)
	logger := xlog.New(output, xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}
