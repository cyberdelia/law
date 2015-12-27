package storage

import (
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/kr/s3/s3util"
)

func init() {
	s3util.DefaultConfig.AccessKey = os.Getenv("AWS_ACCESS_KEY")
	s3util.DefaultConfig.SecretKey = os.Getenv("AWS_SECRET_KEY")
	s3util.DefaultConfig.SecurityToken = os.Getenv("AWS_SECURITY_TOKEN")
}

// S3Storage represents a s3 based file storage.
type S3Storage struct {
	prefix string
}

// NewS3Storage create a new S3Storage base on
// a s3:/// URL.
func NewS3Storage(u *url.URL) *S3Storage {
	u.Scheme = "https"
	return &S3Storage{
		prefix: u.String(),
	}
}

// Create creates a new file based on the given filename.
func (s S3Storage) Create(name string) (io.WriteCloser, error) {
	url := urlJoin(s.prefix, name)
	return s3util.Create(url, nil, nil)
}

// Open opens the given filename.
func (s S3Storage) Open(name string) (io.ReadCloser, error) {
	url := urlJoin(s.prefix, name)
	return s3util.Open(url, nil)
}

// List lists all files present in the file storage after the given prefix.
func (s S3Storage) List(name string) (files []io.ReadCloser, err error) {
	baseurl := urlJoin(s.prefix, name)
	basedir, err := s3util.NewFile(baseurl, nil)
	if err != nil {
		return nil, err
	}
	infos, err := basedir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		file, err := s3util.Open(urlJoin(s.prefix, info.Name()), nil)
		if err != nil {
			return files, nil
		}
		files = append(files, file)
	}
	return files, nil
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
