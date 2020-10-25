package attrs

import (
	"testing"

	"gonum.org/v1/gonum/graph/encoding"
)

func TestAttributes(t *testing.T) {
	a := New()

	exp := 0

	if got := len(a.Attributes()); exp != got {
		t.Errorf("expected %d attributes, got: %d", exp, got)
	}

	keys := a.Keys()
	if count := len(keys); count != exp {
		t.Errorf("expected %d keys, got: %d", exp, count)
	}
}

func TestGetAttribute(t *testing.T) {
	a := New()

	exp := ""
	if val := a.Get("foo"); val != exp {
		t.Errorf("expected: %s, got: %s", exp, val)
	}
}

func TestSetAttribute(t *testing.T) {
	a := New()

	attr := encoding.Attribute{
		Key:   "foo",
		Value: "bar",
	}

	a.Set(attr.Key, attr.Value)

	if val := a.Get(attr.Key); val != attr.Value {
		t.Errorf("expected: %s, got: %s", attr.Value, val)
	}

	exp := 1

	if got := len(a.Attributes()); exp != got {
		t.Errorf("expected %d attributes, got: %d", exp, got)
	}

	keys := a.Keys()

	if count := len(keys); count != exp {
		t.Errorf("expected %d keys, got: %d", exp, count)
	}
}
