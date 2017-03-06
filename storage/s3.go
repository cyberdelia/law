package storage

import (
	"io"
	"net/http"
	"net/url"
	"os"

	s3 "github.com/cyberdelia/s3"
)

// S3Storage represents a s3 based file storage.
type S3Storage struct {
	u      *url.URL
	client *http.Client
}

// NewS3Storage create a new S3Storage base on
// a s3:// URL.
func NewS3Storage(u *url.URL) *S3Storage {
	u.RawQuery = ""
	return &S3Storage{
		u:      u,
		client: s3.DefaultClient,
	}
}

// Create creates a new file based on the given filename.
func (s S3Storage) Create(name string) (io.WriteCloser, error) {
	uri, err := urlJoin(name, s.u)
	if err != nil {
		return nil, err
	}
	return s3.Create(uri, http.Header{
		"x-amz-server-side-encryption": []string{"AES256"},
	}, s.client)
}

// Open opens the given filename.
func (s S3Storage) Open(name string) (io.ReadCloser, error) {
	uri, err := urlJoin(name, s.u)
	if err != nil {
		return nil, err
	}
	r, _, err := s3.Open(uri, s.client)
	return r, err
}

// List lists all files presents in the file storage after the given prefix.
func (s S3Storage) List(name string) (files []io.ReadCloser, err error) {
	uri, err := urlJoin(name, s.u)
	if err != nil {
		return nil, err
	}
	err = s3.Walk(uri, func(path string, info os.FileInfo) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return s3.SkipDir
		}
		uri, err := urlJoin(path, s.u)
		if err != nil {
			return err
		}
		file, _, err := s3.Open(uri, nil)
		if err != nil {
			return err
		}
		files = append(files, file)
		return nil
	}, s.client)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func urlJoin(name string, prefix *url.URL) (string, error) {
	u, err := url.Parse(name)
	if err != nil {
		return "", err
	}
	return prefix.ResolveReference(u).String(), nil
}
