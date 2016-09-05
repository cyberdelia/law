package operator

import (
	"database/sql"
	"fmt"
	"time"

	// load postgres drivers
	_ "github.com/lib/pq"
)

// Database represents the underlying postgres database.
type Database struct {
	*sql.DB
}

// Backup represents a backup.
type Backup struct {
	Name   string
	Offset string
}

// NewDatabase returns a new Database based on the given given
// data source.
func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

// StartBackup starts a new backup.
func (d *Database) StartBackup() (*Backup, error) {
	var name, offset string
	label := fmt.Sprintf("freeze_start_%s", time.Now().UTC().Format(time.RFC3339))
	err := d.QueryRow(`SELECT file_name, lpad(file_offset::text, 8, '0') AS file_offset FROM pg_xlogfile_name_offset(pg_start_backup($1))`, label).Scan(&name, &offset)
	return &Backup{
		Name:   name,
		Offset: offset,
	}, err
}

// StopBackup stops the currently running backup.
func (d *Database) StopBackup() (*Backup, error) {
	var name, offset string
	err := d.QueryRow(`SELECT file_name, lpad(file_offset::text, 8, '0') AS file_offset FROM pg_xlogfile_name_offset(pg_stop_backup())`).Scan(&name, &offset)
	return &Backup{
		Name:   name,
		Offset: offset,
	}, err
}

// Version returns the database version.
func (d *Database) Version() (string, error) {
	var version string
	err := d.QueryRow(`SELECT * FROM version()`).Scan(&version)
	return version, err
}
