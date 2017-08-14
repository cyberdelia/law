package ratio

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"
	"testing/quick"
	"time"
)

func TestWriter(t *testing.T) {
	writer := func(payload []byte) bool {
		buf := new(bytes.Buffer)

		rw := RateLimitedWriter(buf, 2, time.Millisecond)
		defer rw.Close()

		io.Copy(rw, bytes.NewReader(payload))
		return bytes.Equal(buf.Bytes(), payload)
	}
	if err := quick.Check(writer, nil); err != nil {
		t.Error(err)
	}
}

func TestReader(t *testing.T) {
	reader := func(payload []byte) bool {
		buf := new(bytes.Buffer)

		rr := RateLimitedReader(bytes.NewBuffer(payload), 2, time.Millisecond)
		defer rr.Close()

		if _, err := io.Copy(buf, rr); err != nil {
			return false
		}

		return bytes.Equal(buf.Bytes(), payload)
	}
	if err := quick.Check(reader, nil); err != nil {
		t.Error(err)
	}
}

func BenchmarkWriter(b *testing.B) {
	b.StopTimer()
	buf := make([]byte, 2*MB)
	b.SetBytes(2 * MB)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rw := RateLimitedWriter(ioutil.Discard, 1*MB, time.Second)
		rw.Write(buf)
		rw.Close()
	}
}

func BenchmarkReader(b *testing.B) {
	b.StopTimer()
	buf := make([]byte, 2*MB)
	b.SetBytes(2 * MB)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rw := RateLimitedReader(rand.Reader, 1*MB, time.Second)
		rw.Read(buf)
		rw.Close()
	}
}
