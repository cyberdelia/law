package storage

import (
	"io"
	"net/http"
	"net/url"
	"strconv"

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
	u.Scheme = "https"
	u.RawQuery = ""
	return &S3Storage{
		u:      u,
		client: s3.DefaultClient,
	}
}

// Create creates a new file based on the given filename.
func (s S3Storage) Create(name string) (io.WriteCloser, error) {
	u, err := url.Parse(name)
	if err != nil {
		return nil, err
	}
	return s3.Create(s.u.ResolveReference(u).String(), http.Header{
		"x-amz-server-side-encryption": []string{"AES256"},
	}, s.client)
}

// Open opens the given filename.
func (s S3Storage) Open(name string) (io.ReadCloser, error) {
	u, err := url.Parse(name)
	if err != nil {
		return nil, err
	}
	r, _, err := s3.Open(s.u.ResolveReference(u).String(), s.client)
	return r, err
}

func mustParseInt(q url.Values, key string, value int) int {
	v, err := strconv.ParseInt(q.Get(key), 10, 64)
	if err != nil {
		return value
	}
	return int(v)
}

func mustParseInt64(q url.Values, key string, value int64) int64 {
	v, err := strconv.ParseInt(q.Get(key), 10, 64)
	if err != nil {
		return value
	}
	return v
}
