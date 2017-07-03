package operator

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	// load postgres drivers
	_ "github.com/lib/pq"
)

const (
	pgControlData = "pg_controldata"
	pgConfig      = "pg_config"
)

// Database represents the underlying postgres database.
type Database interface {
	StartBackup() (*Backup, error)
	StopBackup() (*Backup, error)
}

type onlineDatabase struct {
	dataSourceName string
}

type offlineDatabase struct {
	dataSourceName string

	backup *Backup
}

// Backup represents a backup.
type Backup struct {
	Name   string
	Offset string
}

// NewDatabase returns a new Database based on the given given
// data source.
func NewDatabase(dsn string) (Database, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "postgres":
		return &onlineDatabase{dataSourceName: dsn}, nil
	case "file":
		return &offlineDatabase{dataSourceName: dsn}, nil
	default:
		return nil, errors.New("unsupported scheme")
	}
}

// StartBackup starts a new backup.
func (on *onlineDatabase) StartBackup() (*Backup, error) {
	db, err := sql.Open("postgres", on.dataSourceName)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var name, offset string
	label := fmt.Sprintf("freeze_start_%s", time.Now().UTC().Format(time.RFC3339))
	if err := db.QueryRow(`SELECT file_name, lpad(file_offset::text, 8, '0') AS file_offset FROM pg_xlogfile_name_offset(pg_start_backup($1))`, label).Scan(&name, &offset); err != nil {
		return nil, err
	}
	return &Backup{
		Name:   name,
		Offset: offset,
	}, nil
}

// StopBackup stops the currently running backup.
func (on *onlineDatabase) StopBackup() (*Backup, error) {
	db, err := sql.Open("postgres", on.dataSourceName)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var name, offset string
	if err := db.QueryRow(`SELECT file_name, lpad(file_offset::text, 8, '0') AS file_offset FROM pg_xlogfile_name_offset(pg_stop_backup())`).Scan(&name, &offset); err != nil {
		return nil, err
	}
	return &Backup{
		Name:   name,
		Offset: offset,
	}, nil
}

func (off *offlineDatabase) StartBackup() (*Backup, error) {
	u, err := url.Parse(off.dataSourceName)
	if err != nil {
		return nil, err
	}
	output, err := exec.Command(pgConfig, "--bindir").Output()
	if err != nil {
		return nil, err
	}
	binary := filepath.Join(strings.TrimSpace(string(output)), pgControlData)
	control, err := exec.Command(binary, u.Path).Output()
	if err != nil {
		return nil, err
	}
	var checkpoint, timeline []byte
	for _, l := range bytes.Split(control, []byte("\n")) {
		f := bytes.Split(l, []byte(":"))
		if len(f) == 2 {
			switch {
			case bytes.Contains(f[0], []byte("Latest checkpoint's REDO location")):
				checkpoint = bytes.TrimSpace(f[1])
			case bytes.Contains(f[0], []byte("Latest checkpoint's TimeLineID")):
				timeline = bytes.TrimSpace(f[1])
			}
		}
	}
	location := bytes.Split(checkpoint, []byte("/"))
	off.backup = &Backup{
		Name:   fmt.Sprintf("%08s%08s%08s", timeline, location[0], location[1][0:2]),
		Offset: fmt.Sprintf("%08s", location[1][0:2]),
	}
	return off.backup, nil
}

func (off *offlineDatabase) StopBackup() (*Backup, error) {
	return off.backup, nil
}
