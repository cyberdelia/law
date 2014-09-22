package operator

import (
	"io"
	"io/ioutil"
	"time"

	"code.google.com/p/go.crypto/openpgp"
	"github.com/cyberdelia/lzo"
	"github.com/cyberdelia/pipeline"
	"github.com/cyberdelia/ratio"
)

func GPGWritePipeline(to []*openpgp.Entity) pipeline.WritePipeline {
	return func(w io.WriteCloser) (io.WriteCloser, error) {
		return openpgp.Encrypt(w, to, nil, nil, nil)
	}
}

func GPGReadPipeline(keyring openpgp.KeyRing) pipeline.ReadPipeline {
	return func(r io.ReadCloser) (io.ReadCloser, error) {
		message, err := openpgp.ReadMessage(r, keyring, nil, nil)
		return ioutil.NopCloser(message.UnverifiedBody), err
	}
}

func LZOWritePipeline(w io.WriteCloser) (io.WriteCloser, error) {
	return lzo.NewWriter(w), nil
}

func LZOReadPipeline(r io.ReadCloser) (io.ReadCloser, error) {
	return lzo.NewReader(r)
}

func RateLimitWritePipeline(size int) pipeline.WritePipeline {
	return func(w io.WriteCloser) (io.WriteCloser, error) {
		return ratio.RateLimitedWriter(w, size, time.Second), nil
	}
}

func RateLimitReadPipeline(size int) pipeline.ReadPipeline {
	return func(r io.ReadCloser) (io.ReadCloser, error) {
		return ratio.RateLimitedReader(r, size, time.Second), nil
	}
}
