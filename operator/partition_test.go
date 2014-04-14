package operator

import (
	"testing"
)

func TestFileString(t *testing.T) {
	file := &File{
		Path: "/tmp/file/path",
		Rel:  "/path",
	}
	if path := file.String(); path != "/tmp/file/path" {
		t.Fatalf("string don't match, wants %s got %s", "/tmp/file/path", path)
	}
}
