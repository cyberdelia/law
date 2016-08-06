package operator

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// File represents an archive file.
type File struct {
	Path     string
	Rel      string
	FileInfo os.FileInfo
}

// String returns the path of a file.
func (f *File) String() string {
	return f.Path
}

// Tape represents an archive.
type Tape []*File

// Copy writes a tar archive of all members.
func (t Tape) Copy(w io.WriteCloser) error {
	archive := tar.NewWriter(w)
	for _, member := range t {
		file, err := os.Open(member.Path)
		if err != nil {
			// File might have been deleted, we can ignore it.
			if isNotExist(err) {
				continue
			}
			return err
		}
		info, err := file.Stat()
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = member.Rel
		if err = archive.WriteHeader(header); err != nil {
			return err
		}
		if !info.IsDir() {
			if _, err = io.Copy(archive, file); err != nil {
				if err != tar.ErrWriteTooLong {
					return err
				}
			}
		}
		if err := file.Close(); err != nil {
			return err
		}
	}
	return archive.Close()
}

func walk(cluster string) (files []*File, err error) {
	err = filepath.Walk(cluster, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// An error occured, stop processing
			return err
		}
		if strings.Contains(path, "pg_xlog") && !info.IsDir() {
			// We don't want to capture WAL files but we want the pg_xlog directory
			return nil
		}
		if info.Name() == "postgresql.conf" || info.Name() == "postmaster.pid" {
			// Ignore configuration and pid files
			return nil
		}
		rel, err := filepath.Rel(cluster, path)
		if err != nil {
			return err
		}
		files = append(files, &File{
			Path:     path,
			Rel:      rel,
			FileInfo: info,
		})
		return nil
	})
	return files, err
}

// Archive creates an archive for the given directory.
func Archive(cluster string) (Tape, error) {
	var archive []*File
	files, err := walk(cluster)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		archive = append(archive, file)
	}
	return archive, nil
}

// Extract extracts the archive to the given directory.
func Extract(cluster string, archive io.ReadCloser) error {
	tr := tar.NewReader(archive)
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				// End of archive
				break
			}
			return err
		}
		filename := filepath.Join(cluster, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			os.MkdirAll(filename, info.Mode())
			continue
		}
		file, err := createFile(filename, info.Mode())
		if err != nil {
			return err
		}
		if _, err = io.Copy(file, tr); err != nil {
			return err
		}
	}
	return nil
}

func createFile(name string, mode os.FileMode) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(name), 0700); err != nil {
		return nil, err
	}
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
}

func isNotExist(err error) bool {
	switch pe := err.(type) {
	case nil:
		return false
	case *os.PathError:
		err = pe.Err
	case *os.LinkError:
		err = pe.Err
	}
	return err == syscall.ENOENT || err == os.ErrNotExist
}