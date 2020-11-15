package owner

import "testing"

func TestStringIn(t *testing.T) {
	if !stringIn("list", []string{"get", "list"}) {
		t.Errorf("expected to provide list")
	}

	if stringIn("list", []string{"get"}) {
		t.Errorf("expected to NOT provide list")
	}
}
