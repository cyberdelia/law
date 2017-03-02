package operator

import (
	"io"
	"time"

	"github.com/cyberdelia/lzo"
	"github.com/cyberdelia/pipeline"
	"github.com/cyberdelia/ratio"
)

// lzoWritePipeline returns a WritePipeline that will compress data.
func lzoWritePipeline(w io.WriteCloser) (io.WriteCloser, error) {
	return lzo.NewWriter(w), nil
}

// lzoReadPipeline returns a ReadPipeline that will decompress data.
func lzoReadPipeline(r io.ReadCloser) (io.ReadCloser, error) {
	return lzo.NewReader(r)
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
