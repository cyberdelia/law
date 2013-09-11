package pipeline

import (
	"github.com/cyberdelia/lzo"
	"io"
)

func LZOWritePipeline(w io.WriteCloser) (io.WriteCloser, error) {
	return lzo.NewWriter(w), nil
}

func LZOReadPipeline(r io.ReadCloser) (io.ReadCloser, error) {
	return lzo.NewReader(r)
}
