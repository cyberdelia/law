package pipeline

import (
	"io"

	"github.com/cyberdelia/lzo"
)

func LZOWritePipeline(w io.WriteCloser) (io.WriteCloser, error) {
	return lzo.NewWriter(w), nil
}

func LZOReadPipeline(r io.ReadCloser) (io.ReadCloser, error) {
	return lzo.NewReader(r)
}
