# pipeline

Pipeline allows you to compose a pipeline of ``io.Reader`` or ``io.Writer``,
allowing you to process data easily. 

## Installation

Download and install :

```
$ go get github.com/cyberdelia/pipeline
```

## Usage

Create your own pipes:

```go
package example

import (
  "io"
  "time"

  "github.com/cyberdelia/pipeline"
  "github.com/cyberdelia/ratio"
)

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

func LZOWritePipeline(w io.WriteCloser) (io.WriteCloser, error) {
  return lzo.NewWriter(w), nil
}

func LZOReadPipeline(r io.ReadCloser) (io.ReadCloser, error) {
  return lzo.NewReader(r)
}
```  

And compose a pipeline:

```go
pipe, err := pipeline.PipeWrite(w, RateLimitWritePipeline(10e6), LZOWritePipeline)
if err != nil {
   ...
}
io.Copy(pipe, file)
```
