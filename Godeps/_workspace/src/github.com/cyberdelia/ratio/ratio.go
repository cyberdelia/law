/*
Rate limited io

To create a rate limited io.Writer:

	// Create a new rate-limited writer.
	rw := ratio.RateLimitedWriter(w, 10e6, time.Second)
	io.Copy(rw, reader)

*/
package ratio

import (
	"io"
	"time"
)

type values struct {
	n   int
	err error
}

type operation struct {
	p       []byte
	values  chan values
	written int
}

type rateLimiter struct {
	stop      chan bool
	worker    chan *operation
	action    func(p []byte) (int, error)
	limit     int
	remaining int
	written   int
	op        *operation
}

func (rl *rateLimiter) reset() {
	rl.remaining = rl.limit
}

func (rl *rateLimiter) run(duration time.Duration) {
	for {
		select {
		case <-rl.stop:
			rl.close()
			return
		case <-time.Tick(duration):
			rl.reset()
			rl.write()
		case op := <-rl.worker:
			rl.op = op
			rl.write()
		}
	}
}

func (rl *rateLimiter) write() {
	if rl.op == nil {
		return
	}

	p := rl.op.p
	if rl.remaining < len(rl.op.p) {
		p = rl.op.p[:rl.remaining]
	}
	rl.remaining -= len(p)

	n, err := rl.action(p)
	rl.op.p = rl.op.p[n:]
	rl.op.written += n

	if len(rl.op.p) == 0 || err != nil {
		rl.op.values <- values{rl.op.written, err}
		rl.op = nil
	}
}

func (rl *rateLimiter) record(p []byte) (int, error) {
	op := &operation{
		p:      p,
		values: make(chan values, 1),
	}
	select {
	case <-rl.stop:
		return 0, io.EOF
	case rl.worker <- op:
		val := <-op.values
		return val.n, val.err
	}
}

func (rl *rateLimiter) close() {
	if rl.op != nil {
		rl.op.values <- values{rl.op.written, io.EOF}
	}
}

// Returns a rate limited io.Writer, allowing to write up to size bytes per duration.
func RateLimitedWriter(w io.Writer, size int, duration time.Duration) io.WriteCloser {
	rl := &rateLimiter{
		action: func(p []byte) (int, error) {
			return w.Write(p)
		},
		limit:     size,
		remaining: size,
		stop:      make(chan bool),
		worker:    make(chan *operation),
	}
	go rl.run(duration)
	return rl
}

// Returns a rate limited io.Reader, allowing to read up to size bytes per duration.
func RateLimitedReader(r io.Reader, size int, duration time.Duration) io.ReadCloser {
	rl := &rateLimiter{
		action: func(p []byte) (int, error) {
			return r.Read(p)
		},
		limit:     size,
		remaining: size,
		stop:      make(chan bool),
		worker:    make(chan *operation),
	}
	go rl.run(duration)
	return rl
}

func (rl *rateLimiter) Read(p []byte) (int, error) {
	return rl.record(p)
}

func (rl *rateLimiter) Write(p []byte) (int, error) {
	return rl.record(p)
}

func (rl *rateLimiter) Close() error {
	select {
	case <-rl.stop:
		return nil
	default:
	}
	close(rl.stop)
	return nil
}
