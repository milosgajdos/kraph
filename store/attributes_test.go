package store

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph/encoding"
)

func TestAttrsKeys(t *testing.T) {
	a := NewAttributes()

	keys := []string{"key1", "key2"}

	for _, key := range keys {
		a.Set(key, "foo")
	}

	attrKeys := a.Keys()

	keyMap := make(map[string]bool)
	attrKeyMap := make(map[string]bool)

	for i := range keys {
		keyMap[keys[i]] = true
	}

	for i := range attrKeys {
		attrKeyMap[attrKeys[i]] = true
	}

	if !reflect.DeepEqual(keyMap, attrKeyMap) {
		t.Errorf("expected keys: %v, got: %v", keys, attrKeys)
	}
}

func TestGetAttribute(t *testing.T) {
	a := NewAttributes()

	exp := ""
	if val := a.Get("foo"); val != exp {
		t.Errorf("expected: %s, got: %s", exp, val)
	}
}

func TestSetAttribute(t *testing.T) {
	a := NewAttributes()

	attr := encoding.Attribute{
		Key:   "foo",
		Value: "bar",
	}

	a.Set(attr.Key, attr.Value)

	if val := a.Get(attr.Key); val != attr.Value {
		t.Errorf("expected: %s, got: %s", attr.Value, val)
	}
}
