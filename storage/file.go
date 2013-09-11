package storage

import (
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type FileStorage struct {
	basedir string
}

func NewFileStorage(u *url.URL) *FileStorage {
	return &FileStorage{
		basedir: u.Path,
	}
}

func (s FileStorage) Open(name string) (io.ReadCloser, error) {
	filename, err := preparePath(s.basedir, name)
	if err != nil {
		return nil, err
	}
	return os.Open(filename)
}

func (s FileStorage) Create(name string) (io.WriteCloser, error) {
	filename, err := preparePath(s.basedir, name)
	if err != nil {
		return nil, err
	}
	return os.Create(filename)
}

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
