package kraph

import "testing"

func TestSupports(t *testing.T) {
	if !provides([]string{"get", "list"}, "list") {
		t.Errorf("expected to provide list")
	}

	if provides([]string{"get"}, "list") {
		t.Errorf("expected to NOT provide list")
	}
}
