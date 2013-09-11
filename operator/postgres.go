package operator

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Database struct {
	*sql.DB
}

type Backup struct {
	Name   string
	Offset string
}

func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

func (d *Database) StartBackup() (backup *Backup, err error) {
	var name, offset string
	label := fmt.Sprintf("freeze_start_%s", time.Now().UTC().Format(time.RFC3339))
	err = d.QueryRow("SELECT file_name, lpad(file_offset::text, 8, '0') AS file_offset FROM pg_xlogfile_name_offset(pg_start_backup($1))", label).Scan(&name, &offset)
	return &Backup{
		Name:   name,
		Offset: offset,
	}, err
}

func (d *Database) StopBackup() (backup *Backup, err error) {
	var name, offset string
	err = d.QueryRow("SELECT file_name, lpad(file_offset::text, 8, '0') AS file_offset FROM pg_xlogfile_name_offset(pg_stop_backup())").Scan(&name, &offset)
	return &Backup{
		Name:   name,
		Offset: offset,
	}, err
}

func (d *Database) Version() (version string, err error) {
	err = d.QueryRow("SELECT * FROM version()").Scan(&version)
	return version, err
}
