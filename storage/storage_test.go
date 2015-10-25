package storage

import (
	"testing"
)

func TestUnsupportedStorage(t *testing.T) {
	_, err := NewStorage("scheme://host/path")
	if err == nil {
		t.Fatal(err)
	}
}

func TestSupportedStorage(t *testing.T) {
	_, err := NewStorage("file:///tmp")
	if err != nil {
		t.Fatal(err)
	}
}
