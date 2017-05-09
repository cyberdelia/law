package operator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestOnlineBackup(t *testing.T) {
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

func TestOfflineBackup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	w, _ := os.Getwd()
	db, err := NewDatabase(fmt.Sprintf("file://%s", filepath.Join(w, "../testdata")))
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
