# xlog

[![Go Reference](https://pkg.go.dev/badge/github.com/goinsane/xlog.svg)](https://pkg.go.dev/github.com/goinsane/xlog)

Package xlog provides leveled and structured logging.

	xlog.Info("this is info log, verbosity 0.")
	xlog.Debug("this is debug log, verbosity 0. it will not be shown.")

	xlog.SetSeverity(xlog.SeverityDebug)
	xlog.Debug("this is debug log, verbosity 0.")
	xlog.V(0).Warning("this is warning log, verbosity 0.")
	xlog.V(1).Warning("this is warning log, verbosity 1. it will not be shown.")

	xlog.SetVerbose(2)
	xlog.V(1).Warning("this is warning log, verbosity 1.")
	xlog.V(2).Error("this is error log, verbosity 2.")
	xlog.V(3).Error("this is error log, verbosity 3. it will not be shown.")

	xlog.WithFields(xlog.Field{Key: "a", Val: "11"}).Info("this is info log, verbosity 0 with fields.")
	xlog.WithFieldKeyVals("x", "1", "y", "2").Info("this is info log, verbosity 0 with fields.")
