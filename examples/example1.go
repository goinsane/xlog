// +build ignore

package main

import (
	"github.com/goinsane/xlog"
)

func main() {
	xlog.SetStackTraceSeverity(xlog.SeverityInfo)
	xlog.SetFlags(xlog.FlagDefault | xlog.FlagShortFile | xlog.FlagShortFunc)
	xlog.WithFieldKeyVals("abc", "def").Info("test\ntest1\n")
}
