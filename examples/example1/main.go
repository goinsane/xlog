package main

import (
	"os"
	"time"

	"github.com/goinsane/xlog"
)

func func1(l *xlog.Logger) {
	l.Info("unknown error\n")
}

func func2(l *xlog.Logger) {
	l.Warning("an error occured: \ntype error")
}

func main() {
	lo := xlog.NewTextLogOutput(os.Stdout, xlog.LogOutputFlagDefault|xlog.LogOutputFlagStackTrace|xlog.LogOutputFlagLongFile|xlog.LogOutputFlagPadding)
	l := xlog.New(lo, xlog.SeverityInfo, 0)
	func1(l)
	go func2(l)
	time.Sleep(time.Second)
}
