package entity

import "testing"

func TestMetadata(t *testing.T) {
	m := make(Metadata)

	if val := m.Get("foo"); val != nil {
		t.Errorf("expected nil, got: %#v", val)
	}

	key, val := "foo", "bar"
	m.Set(key, val)

	if ret := m.Get(key); ret == nil {
		t.Errorf("expected: %s, got: %#v", val, ret)
	}
}
