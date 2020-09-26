package entity

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/kraph/store"
)

func TestEntity(t *testing.T) {
	id, name := "foo", "bar"

	e := New(id, name)

	if id != e.ID() {
		t.Errorf("expected ID %s, got: %s", id, e.ID())
	}

	if name != e.Name() {
		t.Errorf("expected name %s, got: %s", name, e.Name())
	}

	if count := len(e.Attributes()); count != 0 {
		t.Errorf("expected 0 attributes, got: %d", count)
	}

	if meta := e.Metadata(); meta == nil {
		t.Errorf("expected metadata, got: %#v", meta)
	}
}

func TestEntityOpts(t *testing.T) {
	a := store.NewAttributes()
	akey, aval := "foo", "val"
	a.Set(akey, aval)

	m := store.NewMetadata()
	mkey := "foo"
	mval := 5
	m.Set(mkey, mval)

	e := New("foo", "bar", store.Meta(m), store.Attributes(a))

	if count := len(e.Attributes()); count == 0 {
		t.Errorf("expected %d attributes, got: %d", len(e.Attributes()), count)
	}

	if val := e.Attrs().Get(akey); val != aval {
		t.Errorf("expected attribute for key %s: %s, got: %s", akey, aval, val)
	}

	if meta := e.Metadata(); meta == nil {
		t.Errorf("expected metadata, got: %#v", meta)
	}

	if val := e.Metadata().Get(mkey); val.(int) != mval {
		t.Errorf("expected metadata for key %s: %d, got: %d", mkey, mval, val)
	}
}

func TestEntityMetadata(t *testing.T) {
	m := store.NewMetadata()
	mkey := "foo"
	mval := 5
	m.Set(mkey, mval)

	e := New("foo", "bar", store.Meta(m))

	if val := e.Metadata().Get(mkey); !reflect.DeepEqual(val, mval) {
		t.Errorf("expected %v for key %s, got: %v", mval, mkey, val)
	}
}

func TestEntityAttrs(t *testing.T) {
	a := store.NewAttributes()
	akey, aval := "foo", "val"
	a.Set(akey, aval)

	e := New("foo", "bar", store.Attributes(a))

	if count := len(e.Attributes()); count != 1 {
		t.Errorf("expected %d attributes, got: %d", 1, count)
	}
}
