package operator

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/cyberdelia/law/storage"
	"github.com/cyberdelia/pipeline"
)

// Operator contains all operations.
type Operator struct {
	s *storage.Storage
}

// NewOperator creates a new operator.
func NewOperator() (*Operator, error) {
	s, err := storage.NewStorage(os.Getenv("STORAGE_URL"))
	if err != nil {
		return nil, err
	}
	return &Operator{
		s: s,
	}, nil
}

// Unarchive restore the given wal segment to the destination.
func (o *Operator) Unarchive(name string, dest string) error {
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	r, err := o.s.Unarchive(name)
	if err != nil {
		return err
	}
	pipe, err := pipeline.PipeRead(r, lz4ReadPipeline)
	if err != nil {
		return err
	}
	if _, err = io.Copy(file, pipe); err != nil {
		return err
	}
	return pipe.Close()
}

// Archive archives the given wal segment.
func (o *Operator) Archive(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	w, err := o.s.Archive(path.Base(name))
	if err != nil {
		return err
	}
	pipe, err := pipeline.PipeWrite(w, rateLimitWritePipeline(10e6), lz4WritePipeline)
	if err != nil {
		return err
	}
	if _, err = io.Copy(pipe, file); err != nil {
		return err
	}
	return pipe.Close()
}

// Backup backups the given cluster directory.
func (o *Operator) Backup(cluster string) error {
	db, err := NewDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	backup, err := db.StartBackup()
	if err != nil {
		return err
	}
	defer db.StopBackup()
	partitions, err := Partition(cluster)
	if err != nil {
		return err
	}
	for n, part := range partitions {
		w, err := o.s.Backup(backup.Name, backup.Offset, n)
		if err != nil {
			return err
		}
		pipe, err := pipeline.PipeWrite(w, rateLimitWritePipeline(10e6), lz4WritePipeline)
		if err != nil {
			return err
		}
		if err := part.Copy(pipe); err != nil {
			return err
		}
		if err := pipe.Close(); err != nil {
			return err
		}
	}
	backup, err = db.StopBackup()
	if err != nil {
		return err
	}
	if err = db.Close(); err != nil {
		return err
	}
	return nil
}

// Restore a named backup to the given cluster directory.
func (o *Operator) Restore(cluster, name string) error {
	if _, err := os.Stat(path.Join(cluster, "postmaster.pid")); err == nil {
		return errors.New("attempt to overwrite a live data directory")
	}
	readers, err := o.s.Restore(name)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(path.Dir(cluster), 0700); err != nil {
		return err
	}
	for _, r := range readers {
		pipe, err := pipeline.PipeRead(r, lz4ReadPipeline)
		if err != nil {
			return err
		}
		if err = Unite(cluster, pipe); err != nil {
			return err
		}
		if err = pipe.Close(); err != nil {
			return err
		}
	}
	return nil
}
