package metadata

import "testing"

func TestMetadata(t *testing.T) {
	m := New()

	keys := m.Keys()
	if count := len(keys); count > 0 {
		t.Errorf("expected %d keys, got: %d", 0, count)
	}
}

func TestMetadataGet(t *testing.T) {
	m := New()

	if val := m.Get("foo"); val != nil {
		t.Errorf("expected nil, got: %#v", val)
	}
}

func TestMetadataSet(t *testing.T) {
	m := New()

	key, val := "foo", "bar"
	m.Set(key, val)

	if ret := m.Get(key); ret == nil {
		t.Errorf("expected: %s, got: %#v", val, ret)
	}

	keys := m.Keys()
	exp := 1

	if count := len(keys); count != exp {
		t.Errorf("expected %d keys, got: %d", exp, count)
	}
}
