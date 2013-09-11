package storage

import (
	"errors"
	"fmt"
	"io"
	"net/url"
)

const (
	CurrentVersion = "001"
)

var (
	ErrUnsupported = errors.New("unsupported storage")
)

type StorageBackend interface {
	Create(name string) (io.WriteCloser, error)
	Open(name string) (io.ReadCloser, error)
	List(name string) ([]io.ReadCloser, error)
}

type Storage struct {
	b StorageBackend
}

func NewStorage(uri string) (*Storage, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	var backend StorageBackend
	switch u.Scheme {
	case "file":
		backend = NewFileStorage(u)
	case "s3":
		backend = NewS3Storage(u)
	default:
		return nil, ErrUnsupported
	}
	return &Storage{
		b: backend,
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
