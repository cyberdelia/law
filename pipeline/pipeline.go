package pipeline

import (
	"io"
)

type WritePipeline func(w io.WriteCloser) (io.WriteCloser, error)

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
