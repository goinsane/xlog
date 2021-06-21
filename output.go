package xlog

import (
	"bufio"
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Output is an interface for Logger output.
// All of Output implementations must be safe for concurrency.
type Output interface {
	Log(log *Log)
}

type multiOutput []Output

func (m multiOutput) Log(log *Log) {
	for _, o := range m {
		o.Log(log.Duplicate())
	}
}

// MultiOutput creates a output that duplicates its logs to all the provided outputs.
func MultiOutput(outputs ...Output) Output {
	m := make(multiOutput, len(outputs))
	copy(m, outputs)
	return m
}

type asyncOutput struct{ Output }

func (a *asyncOutput) Log(log *Log) {
	go a.Output.Log(log)
}

// AsyncOutput creates a output that doesn't blocks its logs to the provided output.
func AsyncOutput(output Output) Output {
	return &asyncOutput{output}
}

// QueuedOutput is intermediate Output implementation between Logger and given Output.
// QueuedOutput has queueing for unblocking Log() method.
type QueuedOutput struct {
	output      Output
	queue       chan *Log
	ctx         context.Context
	ctxCancel   context.CancelFunc
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
		queue:  make(chan *Log, queueLen),
	}
	q.ctx, q.ctxCancel = context.WithCancel(context.Background())
	go q.worker()
	return
}

// Close closed QueuedOutput. Unused QueuedOutput must be closed for freeing resources.
func (q *QueuedOutput) Close() {
	q.ctxCancel()
}

// Log is implementation of Output.
// If blocking is true, Log method blocks execution until underlying output has finished execution.
// Otherwise, Log method sends log to queue if queue is available. When queue is full, it tries to call OnQueueFull
// function.
func (q *QueuedOutput) Log(log *Log) {
	select {
	case <-q.ctx.Done():
		return
	default:
	}
	if q.blocking != 0 {
		q.queue <- log
		return
	}
	select {
	case q.queue <- log:
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

// TextOutput is an implementation of Output by writing texts to io.Writer w.
type TextOutput struct {
	mu      sync.Mutex
	w       io.Writer
	bw      *bufio.Writer
	flags   Flag
	onError *func(error)
}

// NewTextOutput creates a new TextOutput.
func NewTextOutput(w io.Writer) *TextOutput {
	return &TextOutput{
		w:  w,
		bw: bufio.NewWriter(w),
	}
}

// Log is implementation of Output.
func (t *TextOutput) Log(log *Log) {
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

	if t.flags != 0 {
		log.Flags = t.flags
	}

	var text []byte
	text, err = log.MarshalText()
	if err != nil {
		return
	}

	_, err = t.bw.Write(text)
	if err != nil {
		return
	}
}

// SetWriter sets writer.
func (t *TextOutput) SetWriter(w io.Writer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.w = w
	t.bw = bufio.NewWriter(w)
}

// SetFlags sets flags to override every single Log.Flags if the flags argument different than 0.
// By default, 0.
func (t *TextOutput) SetFlags(flags Flag) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.flags = flags
}

// RegisterOnError registers OnError function to use when error occured.
func (t *TextOutput) RegisterOnError(f func(error)) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&t.onError)), unsafe.Pointer(&f))
}
