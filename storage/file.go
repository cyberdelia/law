package storage

import (
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
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

// List lists all files presents in the file storage after the given prefix.
func (s FileStorage) List(name string) (files []io.ReadCloser, err error) {
	basedir, err := preparePath(s.basedir, name)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		files = append(files, file)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func preparePath(basedir, name string) (string, error) {
	filename := path.Join(basedir, name)
	if err := os.MkdirAll(path.Dir(filename), 0700); err != nil {
		return "", err
	}
	return filename, nil
}
