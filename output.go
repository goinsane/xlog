package xlog

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Output is an interface for Logger output. All of Output implementations must be
// concurrent-safe.
type Output interface {
	Log(msg *Message)
}

type multiOutput []Output

func (m multiOutput) Log(msg *Message) {
	for _, o := range m {
		o.Log(msg.Duplicate())
	}
}

// MultiOutput creates a output that duplicates its logs to all the provided outputs.
func MultiOutput(outputs ...Output) Output {
	m := make(multiOutput, len(outputs))
	copy(m, outputs)
	return m
}

type asyncOutput struct{ Output }

func (a *asyncOutput) Log(msg *Message) {
	go a.Output.Log(msg)
}

// AsyncOutput creates a output that doesn't blocks its logs to the provided output.
func AsyncOutput(output Output) Output {
	return &asyncOutput{output}
}

// QueuedOutput is intermediate Output implementation between Logger and given Output.
// QueuedOutput has queueing for unblocking Log() method.
type QueuedOutput struct {
	ctx         context.Context
	ctxCancel   context.CancelFunc
	output      Output
	queue       chan *Message
	blocking    uint32
	onQueueFull *func()
}

// NewQueuedOutput creates QueuedOutput by given output.
func NewQueuedOutput(output Output, queueLen int) (q *QueuedOutput) {
	if queueLen < 0 {
		queueLen = 0
	}
	q = &QueuedOutput{
		output: output,
		queue:  make(chan *Message, queueLen),
	}
	q.ctx, q.ctxCancel = context.WithCancel(context.Background())
	go q.worker()
	return
}

// Close closed QueuedOutput. Unused QueuedOutput must be closed for freeing resources.
func (q *QueuedOutput) Close() {
	q.ctxCancel()
}

// Log is implementation of Output interface.
// If blocking is true, Log method blocks execution until underlying output has finished execution.
// Otherwise, Log method sends log to queue if queue is available. When queue is full, it tries to call OnQueueFull
// function.
func (q *QueuedOutput) Log(msg *Message) {
	select {
	case <-q.ctx.Done():
		return
	default:
	}
	if q.blocking != 0 {
		q.queue <- msg
		return
	}
	select {
	case q.queue <- msg:
	default:
		if q.onQueueFull != nil && *q.onQueueFull != nil {
			(*q.onQueueFull)()
		}
	}
}

// SetBlocking sets QueuedOutput behavior when queue is full.
func (q *QueuedOutput) SetBlocking(blocking bool) {
	var b uint32
	if blocking {
		b = 1
	}
	atomic.StoreUint32(&q.blocking, b)
}

// RegisterOnQueueFull registers OnQueueFull function to use when queue is full.
func (q *QueuedOutput) RegisterOnQueueFull(f func()) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&q.onQueueFull)), unsafe.Pointer(&f))
}

// WaitForEmpty waits until queue is empty by given context.
func (q *QueuedOutput) WaitForEmpty(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
			if len(q.queue) == 0 {
				return nil
			}
		}
	}
}

func (q *QueuedOutput) worker() {
	for done := false; !done; {
		select {
		case <-q.ctx.Done():
			done = true
		case msg := <-q.queue:
			if q.output != nil {
				q.output.Log(msg)
			}
		}
	}
}

// OutputFlag is type of output flag.
type OutputFlag int

const (
	// OutputFlagDate prints the date in the local time zone: 2009/01/23
	OutputFlagDate = OutputFlag(1 << iota)

	// OutputFlagTime prints the time in the local time zone: 01:23:23
	OutputFlagTime

	// OutputFlagMicroseconds prints microsecond resolution: 01:23:23.123123
	OutputFlagMicroseconds

	// OutputFlagUTC uses UTC rather than the local time zone
	OutputFlagUTC

	// OutputFlagSeverity prints severity level
	OutputFlagSeverity

	// OutputFlagPadding prints padding with multiple lines
	OutputFlagPadding

	// OutputFlagLongFunc prints full package name and function name: a/b/c/d.Func1()
	OutputFlagLongFunc

	// OutputFlagShortFunc prints final package name and function name: d.Func1()
	OutputFlagShortFunc

	// OutputFlagLongFile prints full file name and line number: a/b/c/d.go:23
	OutputFlagLongFile

	// OutputFlagShortFile prints final file name element and line number: d.go:23
	OutputFlagShortFile

	// OutputFlagFields prints fields
	OutputFlagFields

	// OutputFlagStackTrace prints stack trace
	OutputFlagStackTrace

	// OutputFlagDefault holds initial values for the default logger
	OutputFlagDefault = OutputFlagDate | OutputFlagTime | OutputFlagSeverity | OutputFlagFields | OutputFlagStackTrace
)

// TextOutput is an implementation of Output by writing texts to io.Writer w.
type TextOutput struct {
	mu      sync.Mutex
	w       io.Writer
	bw      *bufio.Writer
	flags   OutputFlag
	padding string
	onError *func(error)
}

// NewTextOutput creates a new TextOutput.
func NewTextOutput(w io.Writer, flags OutputFlag) *TextOutput {
	return &TextOutput{
		w:     w,
		bw:    bufio.NewWriter(w),
		flags: flags,
	}
}

// Log implementes Output.Log
func (t *TextOutput) Log(msg *Message) {
	var err error
	defer func() {
		if err == nil || t.onError == nil || *t.onError == nil {
			return
		}
		(*t.onError)(err)
	}()

	t.mu.Lock()
	defer t.mu.Unlock()

	defer func() {
		e := t.bw.Flush()
		if err == nil {
			err = e
		}
	}()

	buf := make([]byte, 128)

	buf = buf[:0]
	if t.flags&(OutputFlagDate|OutputFlagTime|OutputFlagMicroseconds) != 0 {
		if t.flags&OutputFlagUTC != 0 {
			msg.Time = msg.Time.UTC()
		}
		if t.flags&OutputFlagDate != 0 {
			year, month, day := msg.Time.Date()
			itoa(&buf, year, 4)
			buf = append(buf, '/')
			itoa(&buf, int(month), 2)
			buf = append(buf, '/')
			itoa(&buf, day, 2)
			buf = append(buf, ' ')
		}
		if t.flags&(OutputFlagTime|OutputFlagMicroseconds) != 0 {
			hour, min, sec := msg.Time.Clock()
			itoa(&buf, hour, 2)
			buf = append(buf, ':')
			itoa(&buf, min, 2)
			buf = append(buf, ':')
			itoa(&buf, sec, 2)
			if t.flags&OutputFlagMicroseconds != 0 {
				buf = append(buf, '.')
				itoa(&buf, msg.Time.Nanosecond()/1e3, 6)
			}
			buf = append(buf, ' ')
		}
	}

	if t.flags&OutputFlagSeverity != 0 {
		buf = append(buf, msg.Severity.String()...)
		buf = append(buf, " - "...)
	}

	_, err = t.bw.Write(buf)
	if err != nil {
		return
	}

	padding := ""
	if t.flags&OutputFlagPadding != 0 {
		paddingLen := len(buf)
		padding = t.padding
		if padding == "" {
			padding = strings.Repeat(" ", paddingLen)
		}
	}

	if t.flags&(OutputFlagLongFunc|OutputFlagShortFunc) != 0 {
		buf = buf[:0]
		fn := "???"
		if msg.StackCaller.Function != "" {
			fn = msg.StackCaller.Function
		}
		if t.flags&OutputFlagShortFunc != 0 {
			fn = trimDirs(fn)
		}
		buf = append(buf, fn...)
		buf = append(buf, "()"...)
		buf = append(buf, " - "...)
		_, err = t.bw.Write(buf)
		if err != nil {
			return
		}
	}

	if t.flags&(OutputFlagLongFile|OutputFlagShortFile) != 0 {
		buf = buf[:0]
		file, line := "???", 0
		if msg.StackCaller.File != "" {
			file = msg.StackCaller.File
		}
		if msg.StackCaller.Line > 0 {
			line = msg.StackCaller.Line
		}
		if t.flags&OutputFlagShortFile != 0 {
			file = trimDirs(file)
		}
		buf = append(buf, file...)
		buf = append(buf, ':')
		itoa(&buf, line, -1)
		buf = append(buf, " - "...)
		_, err = t.bw.Write(buf)
		if err != nil {
			return
		}
	}

	for i := 0; len(msg.Msg) > 0; i++ {
		if i > 0 {
			_, err = t.bw.WriteString(padding)
			if err != nil {
				return
			}
		}
		idx := bytes.IndexByte(msg.Msg, '\n')
		if idx < 0 {
			msg.Msg = append(msg.Msg, '\n')
			idx = len(msg.Msg)
		} else {
			idx++
		}
		_, err = t.bw.Write(msg.Msg[:idx])
		if err != nil {
			return
		}
		msg.Msg = msg.Msg[idx:]
	}

	if t.flags&OutputFlagFields != 0 && len(msg.Fields) > 0 {
		sort.Sort(msg.Fields)
		buf = buf[:0]
		//buf = append(buf, "\tFields: "...)
		buf = append(buf, "\t"...)
		for _, f := range msg.Fields {
			buf = append(buf, fmt.Sprintf("%s=%q ", f.Key, fmt.Sprintf("%v", f.Val))...)
		}
		buf = append(buf[:len(buf)-1], '\n')
		_, err = t.bw.Write(buf)
		if err != nil {
			return
		}
	}

	if t.flags&OutputFlagStackTrace != 0 && msg.StackTrace != nil {
		buf = buf[:0]
		buf = append(buf, msg.StackTrace.String()...)
		_, err = t.bw.Write(buf)
		if err != nil {
			return
		}
	}
}

// SetWriter sets output writer.
func (t *TextOutput) SetWriter(w io.Writer) {
	t.mu.Lock()
	t.w = w
	t.bw = bufio.NewWriter(w)
	t.mu.Unlock()
}

// SetFlags sets output flags.
func (t *TextOutput) SetFlags(flags OutputFlag) {
	t.mu.Lock()
	t.flags = flags
	t.mu.Unlock()
}

// SetPadding sets custom padding. If padding is empty-string, padding is filled by first line of log.
func (t *TextOutput) SetPadding(padding string) {
	t.mu.Lock()
	t.padding = padding
	t.mu.Unlock()
}

// RegisterOnError registers OnError function to use when error occured.
func (t *TextOutput) RegisterOnError(f func(error)) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&t.onError)), unsafe.Pointer(&f))
}
