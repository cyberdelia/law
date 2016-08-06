package storage

import (
	"io"
	"net/url"
	"os"
	"path"
)

// FileStorage represents a directory based file storage.
type FileStorage struct {
	basedir string
}

// NewFileStorage creates a new FileStorage based on the given
// file:/// URL.
func NewFileStorage(u *url.URL) *FileStorage {
	return &FileStorage{
		basedir: u.Path,
	}
}

// Open opens the given filename.
func (s FileStorage) Open(name string) (io.ReadCloser, error) {
	filename, err := preparePath(s.basedir, name)
	if err != nil {
		return nil, err
	}
	return os.Open(filename)
}

// Create creates a new file based on the given filename.
func (s FileStorage) Create(name string) (io.WriteCloser, error) {
	filename, err := preparePath(s.basedir, name)
	if err != nil {
		return nil, err
	}
	return os.Create(filename)
}

func preparePath(basedir, name string) (string, error) {
	filename := path.Join(basedir, name)
	if err := os.MkdirAll(path.Dir(filename), 0700); err != nil {
		return "", err
	}
	return filename, nil
}
