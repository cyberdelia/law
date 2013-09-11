package pipeline

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"testing"
)

func encodePipeline(w io.WriteCloser) (io.WriteCloser, error) {
	return base64.NewEncoder(base64.StdEncoding, w), nil
}

func decodePipeline(r io.ReadCloser) (io.ReadCloser, error) {
	return ioutil.NopCloser(base64.NewDecoder(base64.StdEncoding, r)), nil
}

type buffer struct {
	bytes.Buffer
}

func (b *buffer) Close() error {
	return nil
}

func TestPipeWriter(t *testing.T) {
	input := bytes.NewBufferString("aloha")
	output := new(buffer)
	w, err := PipeWrite(output, encodePipeline)
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(w, input)
	w.Close()
	if output.String() != "YWxvaGE=" {
		t.Fatal()
	}
}

func TestPipeRead(t *testing.T) {
	input := new(buffer)
	input.WriteString("YWxvaGE=")
	output := new(bytes.Buffer)
	r, err := PipeRead(input, decodePipeline)
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(output, r)
	r.Close()
	if output.String() != "aloha" {
		t.Fatal()
	}
}
