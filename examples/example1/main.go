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
	o := xlog.NewTextOutput(os.Stdout, xlog.OutputFlagDefault|xlog.OutputFlagStackTrace|xlog.OutputFlagLongFile|xlog.OutputFlagPadding)
	l := xlog.New(o, xlog.SeverityInfo, 0)
	func1(l)
	go func2(l)
	time.Sleep(time.Second)
}
