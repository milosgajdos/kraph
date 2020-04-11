package entity

import (
	"testing"

	"gonum.org/v1/gonum/graph/encoding"
)

func TestAttributes(t *testing.T) {
	a := make(Attributes)

	exp := 0
	if got := len(a.Attributes()); exp != got {
		t.Errorf("expected %d attributes, got: %d", exp, got)
	}

	if got := len(a.DOTAttributes()); exp != got {
		t.Errorf("expected %d DOTattributes, got: %d", exp, got)
	}
}

func TestGetAttribute(t *testing.T) {
	a := make(Attributes)

	exp := ""
	if val := a.Get("foo"); val != exp {
		t.Errorf("expected: %s, got: %s", exp, val)
	}
}

func TestSetAttribute(t *testing.T) {
	a := make(Attributes)

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

	if got := len(a.DOTAttributes()); exp != got {
		t.Errorf("expected %d DOTattributes, got: %d", exp, got)
	}
}
