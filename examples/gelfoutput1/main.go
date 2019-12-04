// +build ignore

package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/goinsane/xlog"
	"github.com/goinsane/xlog/gelfoutput"
)

func main() {
	var err error
	output, err := gelfoutput.NewGelfOutput(gelfoutput.GelfWriterTypeTCP, "127.0.0.1:12201", 10, gelfoutput.GelfOptions{})
	if err != nil {
		xlog.Fatal(err)
	}
	defer output.Close()
	output.RegisterOnQueueFull(func() {
		xlog.Error("queue full")
	})
	logger := xlog.New(output, xlog.SeverityInfo, 0)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	for i := 0; ; i++ {
		select {
		case <-sigCh:
			xlog.Info("terminating")
			ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer ctxCancel()
			if err := output.WaitForIdle(ctx); err != nil {
				xlog.Error(err)
			}
			xlog.Info("terminated")
			return
		default:
			logger.WithFieldKeyVals("key1", "val1").Info("test")
			if i > 0 && i%100 == 0 {
				xlog.Infof("sent %d logs", i)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
