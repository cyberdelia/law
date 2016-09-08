package operator

import (
	"compress/gzip"
	"io"
	"time"

	"github.com/cyberdelia/pipeline"
	"github.com/cyberdelia/ratio"
)

// gzipWritePipeline returns a WritePipeline that compress data using GZIP.
func gzipWritePipeline(w io.WriteCloser) (io.WriteCloser, error) {
	return gzip.NewWriter(w), nil
}

// gzipReadPipeline returns a ReadPipeline that decompress data using GZIP.
func gzipReadPipeline(r io.ReadCloser) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

// rateLimitWritePipeline returns a WritePipeline that will rate-limit write I/O.
func rateLimitWritePipeline(size int) pipeline.WritePipeline {
	return func(w io.WriteCloser) (io.WriteCloser, error) {
		if size > 0 {
			return ratio.RateLimitedWriter(w, size, time.Second), nil
		}
		return w, nil
	}
}
