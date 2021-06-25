// Package gelfoutput provides GELF output implementation of xlog.Output.
package gelfoutput

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/goinsane/erf"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"

	"github.com/goinsane/xlog"
)

// GelfOutput implements xlog.Output for GELF output.
type GelfOutput struct {
	opts      Options
	ctx       context.Context
	ctxCancel context.CancelFunc
	mu        sync.Mutex
	writer    gelf.Writer
}

// New creates a new GelfOutput.
func New(opts Options) (g *GelfOutput, err error) {
	g = &GelfOutput{
		opts: opts,
	}
	if g.opts.Host == "" {
		h, e := os.Hostname()
		if e != nil {
			return nil, erf.Errorf("unable to get hostname: %w", e)
		}
		g.opts.Host = h
	}
	g.ctx, g.ctxCancel = context.WithCancel(context.Background())
	return g, nil
}

// Close closes GelfOutput. Unused GelfOutput must be closed for freeing resources.
func (g *GelfOutput) Close() (err error) {
	g.ctxCancel()
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.writer != nil {
		if e := g.writer.Close(); e != nil {
			if err == nil {
				err = erf.Errorf("unable to close writer: %w", e)
			}
		} else {
			g.writer = nil
		}
	}
	return
}

// Log is implementation of xlog.Output.
func (g *GelfOutput) Log(log *xlog.Log) {
	if g.ctx.Err() != nil {
		return
	}
	level := int32(gelf.LOG_EMERG)
	switch log.Severity {
	case xlog.SeverityFatal:
		level = gelf.LOG_CRIT
	case xlog.SeverityError:
		level = gelf.LOG_ERR
	case xlog.SeverityWarning:
		level = gelf.LOG_WARNING
	case xlog.SeverityInfo:
		level = gelf.LOG_INFO
	case xlog.SeverityDebug:
		level = gelf.LOG_DEBUG
	}
	msg := &gelf.Message{
		Version:  "1.1",
		Host:     g.opts.Host,
		Short:    string(log.Message),
		Full:     "",
		TimeUnix: float64(log.Time.UnixNano()) / float64(time.Second),
		Level:    level,
		Facility: g.opts.Facility,
		Extra:    make(map[string]interface{}),
	}
	msg.Extra["severity"] = fmt.Sprintf("%s", log.Severity)
	msg.Extra["verbosity"] = int(log.Verbosity)
	if log.Error != nil {
		format := "%v"
		if _, ok := log.Error.(*erf.Erf); ok {
			format = "%+v"
		}
		msg.Extra["error"] = fmt.Sprintf(format, log.Error)
	}
	msg.Extra["file"] = log.StackCaller.File
	msg.Extra["line"] = log.StackCaller.Line
	msg.Extra["func"] = log.StackCaller.Function
	if log.StackTrace != nil {
		msg.Extra["stack_trace"] = fmt.Sprintf("%+s", log.StackTrace)
	}
	for i := range log.Fields {
		field := &log.Fields[i]
		msg.Extra[fmt.Sprintf("%10.0d_%s", i, field.Key)] = field.Value
	}
	g.writeMessage(msg)
}

func (g *GelfOutput) writeMessage(msg *gelf.Message) {
	var err error
	g.mu.Lock()
	defer g.mu.Unlock()
	for g.ctx.Err() == nil {
		if g.writer == nil {
			if !g.opts.UseTCP {
				var w *gelf.UDPWriter
				w, err = gelf.NewUDPWriter(g.opts.Address)
				if err != nil {
					time.Sleep(250 * time.Millisecond)
					continue
				}
				g.writer = w
			} else {
				var w *gelf.TCPWriter
				w, err = gelf.NewTCPWriter(g.opts.Address)
				if err != nil {
					time.Sleep(250 * time.Millisecond)
					continue
				}
				w.MaxReconnect = 0
				w.ReconnectDelay = 1
				g.writer = w
			}
		}
		if err = g.writer.WriteMessage(msg); err != nil {
			_ = g.writer.Close()
			g.writer = nil
			time.Sleep(250 * time.Millisecond)
			continue
		}
		return
	}
}
