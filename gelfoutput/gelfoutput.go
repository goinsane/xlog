package gelfoutput

import (
	"context"
	"github.com/goinsane/xlog"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"os"
	"sync"
	"time"
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
	mu         sync.Mutex
	ctx        context.Context
	ctxCancel  context.CancelFunc
	writerType GelfWriterType
	addr       string
	opts       GelfOptions
	writer     gelfWriter
}

func NewGelfOutput(writerType GelfWriterType, addr string, opts GelfOptions) (o *GelfOutput, err error) {
	if !writerType.IsValid() {
		err = ErrUnknownGelfWriterType
		return
	}
	o = &GelfOutput{
		writerType: writerType,
		addr:       addr,
		opts:       opts,
	}
	if o.opts.Host == "" {
		osHostname, _ := os.Hostname()
		o.opts.Host = osHostname
	}
	o.ctx, o.ctxCancel = context.WithCancel(context.Background())
	return
}

func (o *GelfOutput) Close() {
	o.ctxCancel()
	o.mu.Lock()
	if o.writer != nil {
		o.writer.Close()
		o.writer = nil
	}
	o.mu.Unlock()
}

func (o *GelfOutput) Log(msg []byte, severity xlog.Severity, verbose xlog.Verbose, tm time.Time, fields xlog.Fields, callers xlog.Callers) {
	select {
	case <-o.ctx.Done():
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
		Host:     o.opts.Host,
		Short:    string(msg),
		Full:     "",
		TimeUnix: float64(tm2.UnixNano()) / float64(time.Second),
		Level:    level,
		Facility: o.opts.Facility,
		Extra:    make(map[string]interface{}),
	}
	m.Extra["severity"] = severity.String()
	for i := range fields {
		field := &fields[i]
		key := "_" + field.Key
		/*if _, ok := m.Extra[key]; ok {
			continue
		}*/
		m.Extra[key] = field.Val
	}
	o.mu.Lock()
	o.writeMessage(m)
	o.mu.Unlock()
}

func (o *GelfOutput) writeMessage(m *gelf.Message) {
	for {
		select {
		case <-o.ctx.Done():
			return
		default:
		}
		var e error
		if o.writer == nil {
			o.writer, e = newGelfWriter(o.writerType, o.addr)
			if e != nil {
				o.writer = nil
				time.Sleep(250 * time.Millisecond)
				continue
			}
		}
		e = o.writer.WriteMessage(m)
		if e != nil {
			o.writer.Close()
			o.writer = nil
			time.Sleep(250 * time.Millisecond)
			continue
		}
		return
	}
}
