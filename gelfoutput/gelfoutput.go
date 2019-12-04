package gelfoutput

import (
	"context"
	"os"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/goinsane/xlog"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

type GelfWriterType int

func (t GelfWriterType) IsValid() bool {
	return t >= GelfWriterTypeUDP && t <= GelfWriterTypeTCP
}

const (
	GelfWriterTypeUDP = iota
	GelfWriterTypeTCP
)

type gelfWriter interface {
	Write(p []byte) (n int, err error)
	WriteMessage(m *gelf.Message) (err error)
	Close() error
}

func newGelfWriter(writerType GelfWriterType, addr string) (w gelfWriter, err error) {
	switch writerType {
	case GelfWriterTypeUDP:
		var gw *gelf.UDPWriter
		gw, err = gelf.NewUDPWriter(addr)
		if err != nil {
			return
		}
		w = gw
	case GelfWriterTypeTCP:
		var gw *gelf.TCPWriter
		gw, err = gelf.NewTCPWriter(addr)
		if err != nil {
			return
		}
		gw.MaxReconnect = 0
		w = gw
	default:
		err = ErrUnknownGelfWriterType
		return
	}
	return
}

type GelfOptions struct {
	Host     string
	Facility string
}

type GelfOutput struct {
	ctx         context.Context
	ctxCancel   context.CancelFunc
	writerType  GelfWriterType
	addr        string
	queue       chan *gelf.Message
	opts        GelfOptions
	w           gelfWriter
	onQueueFull *func()
}

func NewGelfOutput(writerType GelfWriterType, addr string, queueLen int, opts GelfOptions) (g *GelfOutput, err error) {
	if !writerType.IsValid() {
		err = ErrUnknownGelfWriterType
		return
	}
	if queueLen <= 0 {
		queueLen = 1
	}
	g = &GelfOutput{
		writerType: writerType,
		addr:       addr,
		queue:      make(chan *gelf.Message, queueLen),
		opts:       opts,
	}
	if g.opts.Host == "" {
		osHostname, _ := os.Hostname()
		g.opts.Host = osHostname
	}
	g.ctx, g.ctxCancel = context.WithCancel(context.Background())
	go g.worker()
	return
}

func (g *GelfOutput) Close() {
	g.ctxCancel()
}

func (g *GelfOutput) Log(msg []byte, severity xlog.Severity, verbose xlog.Verbose, tm time.Time, fields xlog.Fields, callers xlog.Callers) {
	select {
	case <-g.ctx.Done():
		return
	default:
	}
	tm2 := tm
	if tm2.IsZero() {
		tm2 = time.Now()
	}
	level := int32(gelf.LOG_EMERG)
	switch severity {
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
	m := &gelf.Message{
		Version:  "1.1",
		Host:     g.opts.Host,
		Short:    string(msg),
		Full:     "",
		TimeUnix: float64(tm2.UnixNano()) / float64(time.Second),
		Level:    level,
		Facility: g.opts.Facility,
		Extra:    make(map[string]interface{}),
	}
	for i := range fields {
		field := &fields[i]
		key := "_" + field.Key
		/*if _, ok := m.Extra[key]; ok {
			continue
		}*/
		m.Extra[key] = field.Val
	}
	select {
	case g.queue <- m:
	default:
		if g.onQueueFull != nil && *g.onQueueFull != nil {
			(*g.onQueueFull)()
		}
	}
}

func (g *GelfOutput) RegisterOnQueueFull(f func()) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&g.onQueueFull)), unsafe.Pointer(&f))
}

func (g *GelfOutput) WaitForIdle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
			if len(g.queue) == 0 {
				return nil
			}
		}
	}
}

func (g *GelfOutput) writeMessage(m *gelf.Message) {
	for {
		select {
		case <-g.ctx.Done():
			return
		default:
		}
		var e error
		if g.w == nil {
			time.Sleep(1 * time.Second)
			g.w, e = newGelfWriter(g.writerType, g.addr)
			if e != nil {
				continue
			}
		}
		e = g.w.WriteMessage(m)
		if e == nil {
			break
		}
		g.w.Close()
		g.w = nil
	}
}

func (g *GelfOutput) worker() {
	g.w, _ = newGelfWriter(g.writerType, g.addr)
	for done := false; !done; {
		select {
		case <-g.ctx.Done():
			done = true
		case m := <-g.queue:
			g.writeMessage(m)
		}
	}
	if g.w != nil {
		g.w.Close()
		g.w = nil
	}
}
