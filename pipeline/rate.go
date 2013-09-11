package pipeline

import (
	"github.com/cyberdelia/ratio"
	"io"
	"time"
)

func RateLimitWritePipeline(size int) WritePipeline {
	return func(w io.WriteCloser) (io.WriteCloser, error) {
		return ratio.RateLimitedWriter(w, size, time.Second), nil
	}
}

func RateLimitReadPipeline(size int) ReadPipeline {
	return func(r io.ReadCloser) (io.ReadCloser, error) {
		return ratio.RateLimitedReader(r, size, time.Second), nil
	}
}
