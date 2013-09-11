package operator

import (
	"archive/tar"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	MaxPartitionSize    = 1610612736
	MaxPartitionMembers = int(MaxPartitionSize / 262144)
)

var (
	ErrMemberTooBig = errors.New("file too big for tar partition")
)

type File struct {
	Path     string
	Rel      string
	FileInfo os.FileInfo
}

func (f *File) String() string {
	return f.Path
}

type Archive struct {
	members []*File
}

func (a *Archive) Copy(w io.WriteCloser) error {
	archive := tar.NewWriter(w)
	for _, member := range a.members {
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

func walker(root string, files chan *File) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
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
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files <- &File{
			Path:     path,
			Rel:      rel,
			FileInfo: info,
		}
		return nil
	}
}

func walk(cluster string) chan *File {
	files := make(chan *File)
	go func() {
		err := filepath.Walk(cluster, walker(cluster, files))
		if err != nil {
			log.Println(err)
		}
		close(files)
	}()
	return files
}

func Partition(cluster string) chan *Archive {
	archives := make(chan *Archive)
	go func() {
		var size int64
		var members []*File
		for file := range walk(cluster) {
			if file.FileInfo.Size() > MaxPartitionSize {
				// File is bigger than the max size of partition
				log.Println(ErrMemberTooBig)
			}
			if (size+file.FileInfo.Size() >= MaxPartitionSize) ||
				(len(members) >= MaxPartitionMembers) {
				archives <- &Archive{members}
				members = make([]*File, 0)
				size = 0
			}
			members = append(members, file)
			size += file.FileInfo.Size()
		}
		archives <- &Archive{members}
		close(archives)
	}()
	return archives
}

func Unite(cluster string, partition io.ReadCloser) error {
	archive := tar.NewReader(partition)
	for {
		header, err := archive.Next()
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
		if _, err = io.Copy(file, archive); err != nil {
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
