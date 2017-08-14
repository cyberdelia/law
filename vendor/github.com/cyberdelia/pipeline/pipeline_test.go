// Pipeline allows you to compose read and write pipeline from
// io.Reader and io.Writer.
package pipeline

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"sync/atomic"
	"testing"
)

func writePipeline(count *int64) WritePipeline {
	return func(w io.WriteCloser) (io.WriteCloser, error) {
		atomic.AddInt64(count, 1)
		return w, nil
	}
}

func readPipeline(count *int64) ReadPipeline {
	return func(r io.ReadCloser) (io.ReadCloser, error) {
		atomic.AddInt64(count, 1)
		return r, nil
	}
}

type nopCloser struct {
	io.Writer
}

func (nc nopCloser) Close() error { return nil }

func NopCloser(w io.Writer) io.WriteCloser {
	return nopCloser{w}
}

func TestPipeWriter(t *testing.T) {
	var output bytes.Buffer
	var count int64
	w, err := PipeWrite(NopCloser(&output), writePipeline(&count), writePipeline(&count))
	if err != nil {
		t.Error(err)
	}
	defer w.Close()
	if _, err := io.CopyN(w, rand.Reader, 64); err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Errorf("unexpected count, got %d, wants %d", count, 2)
	}
}

func TestPipeRead(t *testing.T) {
	var output bytes.Buffer
	var count int64
	r, err := PipeRead(ioutil.NopCloser(rand.Reader), readPipeline(&count), readPipeline(&count))
	if err != nil {
		t.Error(err)
	}
	defer r.Close()
	if _, err := io.CopyN(&output, r, 64); err != nil {
		t.Error(err)
	}
}
