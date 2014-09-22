package pipeline

import (
	"io"
)

// WritePipeline type is an adapter to allow the use of different
// writer as pipeline steps.
type WritePipeline func(w io.WriteCloser) (io.WriteCloser, error)

// ReadPipeline type is an adapter to allow the use of different
// reader as pipeline steps.
type ReadPipeline func(r io.ReadCloser) (io.ReadCloser, error)

type pipe struct {
	pipes []interface{}
	pipe  interface{}
}

func (pi *pipe) Write(p []byte) (int, error) {
	w := pi.pipe.(io.WriteCloser)
	return w.Write(p)
}

func (pi *pipe) Read(p []byte) (int, error) {
	r := pi.pipe.(io.ReadCloser)
	return r.Read(p)
}

func (pi *pipe) Close() error {
	for i := len(pi.pipes) - 1; i >= 0; i = i - 1 {
		pipe := pi.pipes[i].(io.Closer)
		if err := pipe.Close(); err != nil {
			return err
		}
	}
	return nil
}

// PipeWrite returns a WriteCloser that's the logical concatenation of the provided input
// writer and the given WritePipeline.
func PipeWrite(w io.WriteCloser, pipelines ...WritePipeline) (io.WriteCloser, error) {
	var pipes []interface{}
	pipes = append(pipes, w)
	for _, pipeline := range pipelines {
		w, err := pipeline(w)
		if err != nil {
			return nil, err
		}
		pipes = append(pipes, w)
	}
	return &pipe{
		pipes: pipes,
		pipe:  pipes[len(pipes)-1],
	}, nil
}

// PipeRead returns a ReadCloser that's the logical concatenation of the provided input
// reader and the given ReadPipeline.
func PipeRead(r io.ReadCloser, pipelines ...ReadPipeline) (io.ReadCloser, error) {
	var pipes []interface{}
	pipes = append(pipes, r)
	for _, pipeline := range pipelines {
		r, err := pipeline(r)
		if err != nil {
			return nil, err
		}
		pipes = append(pipes, r)
	}
	return &pipe{
		pipes: pipes,
		pipe:  pipes[len(pipes)-1],
	}, nil
}
