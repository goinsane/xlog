// +build examples

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/goinsane/xlog"
	"github.com/goinsane/xlog/gelfoutput"
)

func main() {
	var address string
	var useTCP bool
	flag.StringVar(&address, "a", "127.0.0.1:12201", "graylog address")
	flag.BoolVar(&useTCP, "t", false, "use tcp")
	flag.Parse()

	xlog.SetOutputWriter(os.Stdout)

	var err error
	gelfOutput, err := gelfoutput.New(gelfoutput.Options{
		Address: address,
		UseTCP:  useTCP,
	})
	if err != nil {
		xlog.Fatal(err)
	}
	defer gelfOutput.Close()

	queuedOutput := xlog.NewQueuedOutput(gelfOutput, 10)
	defer queuedOutput.Close()
	queuedOutput.SetOnQueueFull(func() {
		xlog.Error("queue full")
	})

	logger := xlog.New(queuedOutput, xlog.SeverityInfo, 0)
	logger.SetStackTraceSeverity(xlog.SeverityInfo)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	for i := 0; ; i++ {
		select {
		case <-sigCh:
			xlog.Info("terminating")
			ctx, ctxCancel := context.WithTimeout(context.Background(), 2*time.Second)
			if err := queuedOutput.WaitForEmpty(ctx); err != nil {
				xlog.Error(err)
			}
			ctxCancel()
			xlog.Info("terminated")
			return
		default:
			logger.WithFieldKeyVals("key1", "val1", "key2", "val2", "key1", "val1-2", "key1", "val1-3").Info("test")
			if i > 0 && i%100 == 0 {
				xlog.Infof("sent %d logs", i)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
