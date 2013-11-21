package operator

import (
	"os"
	"testing"
)

func TestNewOperator(t *testing.T) {
	os.Setenv("STORAGE_URL", "file:///tmp")
	if _, err := NewOperator(); err != nil {
		t.Fatal(err)
	}
}
