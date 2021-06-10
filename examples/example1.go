// +build ignore

package main

import (
	"github.com/goinsane/xlog"
)

func main() {
	xlog.SetStackTraceSeverity(xlog.SeverityInfo)
	xlog.Info("test")
}
