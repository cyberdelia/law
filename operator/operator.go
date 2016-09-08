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
func NewOperator(ssn string) (*Operator, error) {
	s, err := storage.NewStorage(ssn)
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
	pipe, err := pipeline.PipeRead(r, gzipReadPipeline)
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
	pipe, err := pipeline.PipeWrite(w, gzipWritePipeline)
	if err != nil {
		return err
	}
	if _, err = io.Copy(pipe, file); err != nil {
		return err
	}
	return pipe.Close()
}

// Backup backups the given cluster directory.
func (o *Operator) Backup(cluster string, rate int) error {
	db, err := NewDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()
	backup, err := db.StartBackup()
	if err != nil {
		return err
	}
	defer db.StopBackup()
	archive, err := Archive(cluster)
	if err != nil {
		return err
	}
	w, err := o.s.Backup(backup.Name, backup.Offset)
	if err != nil {
		return err
	}
	defer w.Close()
	pipe, err := pipeline.PipeWrite(w, rateLimitWritePipeline(rate), gzipWritePipeline)
	if err != nil {
		return err
	}
	defer pipe.Close()
	if err := archive.Copy(pipe); err != nil {
		return err
	}
	return nil
}

// Restore a named backup to the given cluster directory.
func (o *Operator) Restore(cluster, name, offset string) error {
	if _, err := os.Stat(path.Join(cluster, "postmaster.pid")); err == nil {
		return errors.New("attempt to overwrite a live data directory")
	}
	r, err := o.s.Restore(name, offset)
	if err != nil {
		return err
	}
	defer r.Close()
	if err = os.MkdirAll(path.Dir(cluster), 0700); err != nil {
		return err
	}
	pipe, err := pipeline.PipeRead(r, gzipReadPipeline)
	if err != nil {
		return err
	}
	if err = Extract(cluster, pipe); err != nil {
		return err
	}
	return pipe.Close()
}
