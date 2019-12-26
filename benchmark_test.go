package xlog_test

import (
	"io/ioutil"
	"testing"

	"github.com/goinsane/xlog"
)

func BenchmarkLogger(b *testing.B) {
	output := xlog.NewTextOutput(ioutil.Discard, xlog.OutputFlagDefault)
	logger := xlog.New(output, xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}

func BenchmarkLoggerWithStackTrace(b *testing.B) {
	output := xlog.NewTextOutput(ioutil.Discard, xlog.OutputFlagDefault)
	logger := xlog.New(output, xlog.SeverityInfo, 0)
	logger.SetStackTraceSeverity(xlog.SeverityInfo)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark")
	}
}
