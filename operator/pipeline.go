package operator

import (
	"io"
	"io/ioutil"
	"time"

	"github.com/cyberdelia/lzo"
	"github.com/cyberdelia/pipeline"
	"github.com/cyberdelia/ratio"
	"golang.org/x/crypto/openpgp"
)

// GPGWritePipeline returns a WritePipeline that will encrypt data.
func GPGWritePipeline(to []*openpgp.Entity) pipeline.WritePipeline {
	return func(w io.WriteCloser) (io.WriteCloser, error) {
		return openpgp.Encrypt(w, to, nil, nil, nil)
	}
}

// GPGReadPipeline returns a ReadPipeline that will decrypt data.
func GPGReadPipeline(keyring openpgp.KeyRing) pipeline.ReadPipeline {
	return func(r io.ReadCloser) (io.ReadCloser, error) {
		message, err := openpgp.ReadMessage(r, keyring, nil, nil)
		return ioutil.NopCloser(message.UnverifiedBody), err
	}
}

// LZOWritePipeline returns a WritePipeline that will compress data.
func LZOWritePipeline(w io.WriteCloser) (io.WriteCloser, error) {
	return lzo.NewWriter(w), nil
}

// LZOReadPipeline returns a ReadPipeline that will decompress data.
func LZOReadPipeline(r io.ReadCloser) (io.ReadCloser, error) {
	return lzo.NewReader(r)
}

// RateLimitWritePipeline returns a WritePipeline that will rate-limit write I/O.
func RateLimitWritePipeline(size int) pipeline.WritePipeline {
	return func(w io.WriteCloser) (io.WriteCloser, error) {
		return ratio.RateLimitedWriter(w, size, time.Second), nil
	}
}

// RateLimitReadPipeline returns a ReadPipeline that will rate-limit write I/O.
func RateLimitReadPipeline(size int) pipeline.ReadPipeline {
	return func(r io.ReadCloser) (io.ReadCloser, error) {
		return ratio.RateLimitedReader(r, size, time.Second), nil
	}
}
