package storage

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	s3 "github.com/rlmcpherson/s3gof3r"
)

// S3Storage represents a s3 based file storage.
type S3Storage struct {
	prefix string
	bucket *s3.Bucket
	config *s3.Config
}

// NewS3Storage create a new S3Storage base on
// a s3:/// URL.
func NewS3Storage(u *url.URL) *S3Storage {
	k, err := envKeys()
	if err != nil {
		panic(err)
	}
	storage := s3.New(u.Host, k)
	var prefix string
	parts := strings.SplitN(u.Path[1:], "/", 2)
	if len(parts) > 1 {
		prefix = parts[1]
	}
	v := u.Query()
	return &S3Storage{
		prefix: prefix,
		bucket: storage.Bucket(parts[0]),
		config: &s3.Config{
			Concurrency: mustParseInt(v, "concurrency", 10),
			PartSize:    mustParseInt64(v, "part_size", 20000000),
			NTry:        mustParseInt(v, "retry", 10),
			Md5Check:    false,
			Scheme:      "https",
			Client:      s3.ClientWithTimeout(5 * time.Second),
		},
	}
}

// Create creates a new file based on the given filename.
func (s S3Storage) Create(name string) (io.WriteCloser, error) {
	return s.bucket.PutWriter(urlJoin(s.prefix, name), nil, s.config)
}

// Open opens the given filename.
func (s S3Storage) Open(name string) (io.ReadCloser, error) {
	r, _, err := s.bucket.GetReader(urlJoin(s.prefix, name), s.config)
	return r, err
}

func urlJoin(strs ...string) string {
	ss := make([]string, len(strs))
	for i, s := range strs {
		if i == 0 {
			ss[i] = strings.TrimRight(s, "/")
		} else {
			ss[i] = strings.TrimLeft(s, "/")
		}
	}
	return strings.Join(ss, "/")
}

func envKeys() (s3.Keys, error) {
	keys := s3.Keys{
		AccessKey:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:     os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SecurityToken: os.Getenv("AWS_SECURITY_TOKEN"),
	}
	if keys.AccessKey == "" || keys.SecretKey == "" {
		return keys, fmt.Errorf("keys not set in environment: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY")
	}
	return keys, nil
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
