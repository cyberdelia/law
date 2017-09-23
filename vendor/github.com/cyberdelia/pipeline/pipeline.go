package pipeline

import "io"

// WritePipeline type is an adapter to allow the use of different
// writer as pipeline steps.
type WritePipeline func(w io.WriteCloser) (io.WriteCloser, error)

type writePipe struct {
	pipes []io.WriteCloser
	pipe  io.WriteCloser
}

// PipeWrite returns a WriteCloser that's the logical concatenation of the provided input
// writer and the given WritePipeline.
func PipeWrite(w io.WriteCloser, pipelines ...WritePipeline) (io.WriteCloser, error) {
	pipes := []io.WriteCloser{w}
	for _, pipeline := range pipelines {
		w, err := pipeline(pipes[len(pipes)-1].(io.WriteCloser))
		if err != nil {
			return nil, err
		}
		pipes = append(pipes, w)
	}
	return &writePipe{
		pipes: pipes,
		pipe:  pipes[len(pipes)-1],
	}, nil
}

func (wp *writePipe) Write(p []byte) (int, error) {
	return wp.pipe.Write(p)
}

func (wp *writePipe) Close() error {
	for i := len(wp.pipes) - 1; i > 0; i = i - 1 {
		pipe := wp.pipes[i].(io.Closer)
		if err := pipe.Close(); err != nil {
			return err
		}
	}
	return nil
}

// ReadPipeline type is an adapter to allow the use of different
// reader as pipeline steps.
type ReadPipeline func(r io.ReadCloser) (io.ReadCloser, error)

type readPipe struct {
	pipes []io.ReadCloser
	pipe  io.ReadCloser
}

// PipeRead returns a ReadCloser that's the logical concatenation of the provided input
// reader and the given ReadPipeline.
func PipeRead(r io.ReadCloser, pipelines ...ReadPipeline) (io.ReadCloser, error) {
	pipes := []io.ReadCloser{r}
	for _, pipeline := range pipelines {
		r, err := pipeline(pipes[len(pipes)-1].(io.ReadCloser))
		if err != nil {
			return nil, err
		}
		pipes = append(pipes, r)
	}
	return &readPipe{
		pipes: pipes,
		pipe:  pipes[len(pipes)-1],
	}, nil
}

func (rp *readPipe) Read(p []byte) (int, error) {
	return rp.pipe.Read(p)
}

func (rp *readPipe) Close() error {
	for i := len(rp.pipes) - 1; i > 0; i = i - 1 {
		pipe := rp.pipes[i].(io.Closer)
		if err := pipe.Close(); err != nil {
			return err
		}
	}
	return nil
}
