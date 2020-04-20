package entity

import (
	"testing"

	"github.com/milosgajdos/kraph/store"
)

func TestEntity(t *testing.T) {
	e := New()

	if count := len(e.Attrs().Attributes()); count != 0 {
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

	e := New(store.Meta(m), store.EntAttrs(a))

	if count := len(e.Attrs().Attributes()); count == 0 {
		t.Errorf("expected %d attributes, got: %d", len(a.Attributes()), count)
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

func TestEntityAttrs(t *testing.T) {
	a := store.NewAttributes()
	akey, aval := "foo", "val"
	a.Set(akey, aval)

	m := store.NewMetadata()
	mkey := "foo"
	mval := 5
	m.Set(mkey, mval)

	e := New(store.Meta(m), store.EntAttrs(a))

	if count := len(e.Attributes()); count == 0 {
		t.Errorf("expected %d attributes, got: %d", len(a.Attributes()), count)
	}
}
