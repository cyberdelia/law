package operator

import (
	"io"
	"io/ioutil"
	"time"

	"github.com/cyberdelia/pipeline"
	"github.com/cyberdelia/ratio"
	"github.com/pierrec/lz4"
)

// lz4WritePipeline returns a WritePipeline that will compress data.
func lz4WritePipeline(w io.WriteCloser) (io.WriteCloser, error) {
	return lz4.NewWriter(w), nil
}

// lz4ReadPipeline returns a ReadPipeline that will decompress data.
func lz4ReadPipeline(r io.ReadCloser) (io.ReadCloser, error) {
	return ioutil.NopCloser(lz4.NewReader(r)), nil
}

// rateLimitWritePipeline returns a WritePipeline that will rate-limit write I/O.
func rateLimitWritePipeline(size int) pipeline.WritePipeline {
	return func(w io.WriteCloser) (io.WriteCloser, error) {
		return ratio.RateLimitedWriter(w, size, time.Second), nil
	}
}

// rateLimitReadPipeline returns a ReadPipeline that will rate-limit write I/O.
func rateLimitReadPipeline(size int) pipeline.ReadPipeline {
	return func(r io.ReadCloser) (io.ReadCloser, error) {
		return ratio.RateLimitedReader(r, size, time.Second), nil
	}
}
