package storage

import (
	"fmt"
	"io"
	"net/url"
)

// CurrentVersion is a version prefix to be used by storage backends.
const CurrentVersion = "001"

// Backend represents a storage backend.
type Backend interface {
	Create(name string) (io.WriteCloser, error)
	Open(name string) (io.ReadCloser, error)
}

// Storage represents a storage facility.
type Storage struct {
	b Backend
}

// NewStorage create a new storage facility, using
// the appropriate storage backend.
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
		return nil, fmt.Errorf("unsupported storage : %s", u.Scheme)
	}
	return &Storage{
		b: b,
	}, nil
}

// Archive returns a writer to archive the given wal segment.
func (s Storage) Archive(name string) (io.WriteCloser, error) {
	filename := fmt.Sprintf("law_%s/%s.gz", CurrentVersion, name)
	return s.b.Create(filename)
}

// Unarchive returns a reader to restore the given wal segment.
func (s Storage) Unarchive(name string) (io.ReadCloser, error) {
	filename := fmt.Sprintf("law_%s/%s.gz", CurrentVersion, name)
	return s.b.Open(filename)
}

// Backup returns a writer to archive the given backup.
func (s Storage) Backup(name, offset string) (io.WriteCloser, error) {
	filename := fmt.Sprintf("basebackup_%s/%s_%s.tar.gz", CurrentVersion, name, offset)
	return s.b.Create(filename)
}

// Restore returns a reader to restore the given backup.
func (s Storage) Restore(name, offset string) (io.ReadCloser, error) {
	filename := fmt.Sprintf("basebackup_%s/%s_%s.tar.gz", CurrentVersion, name, offset)
	return s.b.Open(filename)
}
