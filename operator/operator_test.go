package operator

import "testing"

func TestNewOperator(t *testing.T) {
	if _, err := NewOperator("file:///tmp"); err != nil {
		t.Fatal(err)
	}
}
