package operator

import (
	"os"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	db, err := NewDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	version, err := db.Version()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(version, "PostgreSQL") {
		t.Error("did not return proper version string")
	}
}

func TestBackup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	db, err := NewDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	start, err := db.StartBackup()
	if err != nil {
		t.Fatal(err)
	}
	stop, err := db.StopBackup()
	if err != nil {
		t.Fatal(err)
	}
	if start.Name != stop.Name {
		t.Error("did not return the same backup name")
	}
}
