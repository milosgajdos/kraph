package gen

import (
	"testing"

	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/uuid"
)

const (
	objPath = "seeds/objects.yaml"
)

func TestObjects(t *testing.T) {
	top, err := NewMockTop(objPath)
	if err != nil {
		t.Errorf("failed to create mock Top: %v", err)
		return
	}

	objects := top.Objects()
	if len(objects) == 0 {
		t.Errorf("no objects found")
	}
}

func TestGetUID(t *testing.T) {
	top, err := NewMockTop(objPath)
	if err != nil {
		t.Errorf("failed to create mock Top: %v", err)
		return
	}

	uids := make([]uuid.UID, len(top.Objects()))
	for i, o := range top.Objects() {
		uids[i] = o.UID()
	}

	for _, uid := range uids {
		q := query.Build().UID(uid, query.UIDEqFunc(uid))

		objects, err := top.Get(q)

		if err != nil {
			t.Errorf("error getting object: %s: %v", uid, err)
			continue
		}

		if len(objects) != 1 {
			t.Errorf("expected 1 %s object, got: %d", uid, len(objects))
			continue
		}

		if objects[0].UID().String() != uid.String() {
			t.Errorf("expected object %s, got: %s", uid, objects[0].UID())
		}
	}
}

func TestTopGet(t *testing.T) {
	top, err := NewMockTop(objPath)
	if err != nil {
		t.Errorf("failed to create mock Top: %v", err)
		return
	}

	q := query.Build().MatchAny()

	objects, err := top.Get(q)
	if err != nil {
		t.Errorf("error getting all objects: %v", err)
	}

	if len(objects) != len(top.Objects()) {
		t.Errorf("expected %d object, got: %d", len(objects), len(top.Objects()))

	}

	namespaces := make([]string, len(top.Objects()))
	kinds := make([]string, len(top.Objects()))
	names := make([]string, len(top.Objects()))

	for i, o := range top.Objects() {
		namespaces[i] = o.Namespace()
		kinds[i] = o.Resource().Kind()
		names[i] = o.Name()
	}

	for _, ns := range namespaces {
		q := query.Build().Namespace(ns, query.StringEqFunc(ns))

		objects, err := top.Get(q)
		if err != nil {
			t.Errorf("error getting namespace %s objects: %v", ns, err)
			continue
		}

		for _, o := range objects {
			if o.Namespace() != ns {
				t.Errorf("expected: namespace %s, got: %s", ns, o.Namespace())
			}
		}

		for _, kind := range kinds {
			q = q.Kind(kind, query.StringEqFunc(kind))

			objects, err = top.Get(q)
			if err != nil {
				t.Errorf("error getting objects: %s/%s: %v", ns, kind, err)
				continue
			}

			for _, o := range objects {
				if o.Namespace() != ns || o.Resource().Kind() != kind {
					t.Errorf("expected: %s/%s, got: %s/%s", ns, kind, o.Namespace(), o.Resource().Kind())
				}
			}

			for _, name := range names {
				q = q.Name(name, query.StringEqFunc(name))

				objects, err = top.Get(q)
				if err != nil {
					t.Errorf("error getting objects: %s/%s/%s: %v", ns, kind, name, err)
					continue
				}

				for _, o := range objects {
					if o.Namespace() != ns || o.Resource().Kind() != kind || o.Name() != name {
						t.Errorf("expected: %s/%s/%s, got: %s/%s/%s", ns, kind, name,
							o.Namespace(), o.Resource().Kind(), o.Name())
					}
				}
			}
		}
	}
}
