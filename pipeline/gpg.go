package pipeline

import (
	"code.google.com/p/go.crypto/openpgp"
	"io"
	"io/ioutil"
)

func GPGWritePipeline(to []*openpgp.Entity) WritePipeline {
	return func(w io.WriteCloser) (io.WriteCloser, error) {
		return openpgp.Encrypt(w, to, nil, nil, nil)
	}
}

func GPGReadPipeline(keyring openpgp.KeyRing) ReadPipeline {
	return func(r io.ReadCloser) (io.ReadCloser, error) {
		message, err := openpgp.ReadMessage(r, keyring, nil, nil)
		return ioutil.NopCloser(message.UnverifiedBody), err
	}
}
