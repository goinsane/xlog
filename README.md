# xlog

[![GoDoc](https://godoc.org/github.com/goinsane/xlog?status.svg)](https://godoc.org/github.com/goinsane/xlog)

Package xlog provides leveled and structured logging.

    xlog.SetSeverity(xlog.SeverityDebug)
    xlog.SetVerbose(2)
    xlog.Info("this is info log, verbosity 0")
    xlog.V(0).Warning("this is warning log, verbosity 0")
    xlog.V(1).Warning("this is warning log, verbosity 1")
    xlog.V(2).Error("this is error log, verbosity 2")
    xlog.V(3).Error("this is error log, verbosity 3. it will not be shown")
    xlog.Debug("this is debug log, verbosity 0")
