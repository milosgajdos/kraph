package kraph

import (
	"testing"

	"gonum.org/v1/gonum/graph/encoding"
)

func TestAttrs(t *testing.T) {
	a := make(Attrs)

	exp := 0
	if got := len(a.Attributes()); exp != got {
		t.Errorf("expected %d attributes, got: %d", exp, got)
	}

	if got := len(a.DOTAttrs()); exp != got {
		t.Errorf("expected %d DOTattributes, got: %d", exp, got)
	}
}

func TestGetAttribute(t *testing.T) {
	a := make(Attrs)

	exp := ""
	if val := a.GetAttribute("foo"); val != exp {
		t.Errorf("expected: %s, got: %s", exp, val)
	}
}

func TestSetAttribute(t *testing.T) {
	a := make(Attrs)

	attr := encoding.Attribute{
		Key:   "foo",
		Value: "bar",
	}

	if err := a.SetAttribute(attr.Key, attr.Value); err != nil {
		t.Fatalf("failed to add attribute: %v", err)
	}

	if val := a.GetAttribute(attr.Key); val != attr.Value {
		t.Errorf("expected: %s, got: %s", attr.Value, val)
	}

	exp := 1

	if got := len(a.Attributes()); exp != got {
		t.Errorf("expected %d attributes, got: %d", exp, got)
	}

	if got := len(a.DOTAttrs()); exp != got {
		t.Errorf("expected %d DOTattributes, got: %d", exp, got)
	}
}
