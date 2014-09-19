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
		t.Fatal("did not return proper version string")
	}
}
