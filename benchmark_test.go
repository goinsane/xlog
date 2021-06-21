package xlog_test

import (
	"io/ioutil"
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
	tm := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithTime(tm)
	}
}

func BenchmarkLogger_WithFieldKeyVals(b *testing.B) {
	logger := xlog.New(xlog.NewTextOutput(ioutil.Discard), xlog.SeverityInfo, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFieldKeyVals("key1", "value1")
	}
}
