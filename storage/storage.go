package storage

import (
	"errors"
	"fmt"
	"io"
	"net/url"
)

const CurrentVersion = "001"

type Backend interface {
	Create(name string) (io.WriteCloser, error)
	Open(name string) (io.ReadCloser, error)
	List(name string) ([]io.ReadCloser, error)
}

type Storage struct {
	b Backend
}

func NewStorage(uri string) (*Storage, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	var b Backend
	switch u.Scheme {
	case "file":
		b = NewFileStorage(u)
	case "s3":
		b = NewS3Storage(u)
	default:
		return nil, errors.New("unsupported storage")
	}
	return &Storage{
		b: b,
	}, nil
}

func (s Storage) Archive(name string) (io.WriteCloser, error) {
	filename := fmt.Sprintf("law_%s/%s.lzo", CurrentVersion, name)
	return s.b.Create(filename)
}

func (s Storage) Unarchive(name string) (io.ReadCloser, error) {
	filename := fmt.Sprintf("law_%s/%s.lzo", CurrentVersion, name)
	return s.b.Open(filename)
}

func (s Storage) Backup(name, offset string, n int) (io.WriteCloser, error) {
	filename := fmt.Sprintf("basebackup_%s/base_%s_%s/part_%d.tar.lzo", CurrentVersion, name, offset, n)
	return s.b.Create(filename)
}

func (s Storage) Restore(name string) ([]io.ReadCloser, error) {
	filename := fmt.Sprintf("basebackup_%s/%s", CurrentVersion, name)
	return s.b.List(filename)
}
