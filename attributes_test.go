package kraph

import (
	"testing"

	"gonum.org/v1/gonum/graph/encoding"
)

func TestAttributes(t *testing.T) {
	var a Attributes

	if len(a.Attributes()) != 0 {
		t.Errorf("expected %d attributes, got: %d", 0, len(a.Attributes()))
	}

	if len(a.DOTAttributes()) != 0 {
		t.Errorf("expected %d attributes, got: %d", 0, len(a.DOTAttributes()))
	}
}

func TestGetSetAttributes(t *testing.T) {
	var a Attributes

	if val := a.Get("foo"); val != "" {
		t.Errorf("expected empty string, got: %s", val)
	}

	attr := encoding.Attribute{
		Key:   "foo",
		Value: "bar",
	}
	if err := a.SetAttribute(attr); err != nil {
		t.Fatalf("failed to set attribute %s: %v", attr.Key, err)
	}

	if val := a.Get(attr.Key); val != attr.Value {
		t.Errorf("expected: %s, got: %s", attr.Value, val)
	}

	attr.Value = "bar2"
	if err := a.SetAttribute(attr); err != nil {
		t.Fatalf("failed to set attribute %s: %v", attr.Key, err)
	}

	if val := a.Get(attr.Key); val != attr.Value {
		t.Errorf("expected: %s, got: %s", attr.Value, val)
	}
}
