package mock

import (
	"strings"
	"testing"

	"github.com/milosgajdos/kraph/query"
)

func TestObjects(t *testing.T) {
	top := NewTop()

	objects := top.Objects()
	if len(objects) == 0 {
		t.Errorf("no objects found")
	}
}

func TestGetUID(t *testing.T) {
	top := NewTop()

	for _, nsKinds := range ObjectData {
		for nsKind, names := range nsKinds {
			nsplit := strings.Split(nsKind, "/")
			ns, kind := nsplit[0], nsplit[1]
			for _, name := range names {
				uid := strings.Join([]string{ns, kind, name}, "/")
				objects, err := top.Get(query.UID(uid))
				if err != nil {
					t.Errorf("error getting object: %s: %v", uid, err)
					continue
				}

				if len(objects) != 1 {
					t.Errorf("expected 1 object, got: %d", len(objects))
					continue
				}

				if objects[0].UID().String() != uid {
					t.Errorf("expected object %s, got: %s", uid, objects[0].UID())
				}
			}
		}
	}
}

func TestGetNsKind(t *testing.T) {
	top := NewTop()

	for _, nsKinds := range ObjectData {
		for nsKind, _ := range nsKinds {
			nsplit := strings.Split(nsKind, "/")
			ns, kind := nsplit[0], nsplit[1]
			objects, err := top.Get(query.Namespace(ns), query.Kind(kind))
			if err != nil {
				t.Errorf("error getting objects: %s/%s: %v", ns, kind, err)
				continue
			}

			for _, object := range objects {
				if object.Namespace() != ns || object.Kind() != kind {
					t.Errorf("expected: %s/%s, got: %s/%s", ns, kind, object.Namespace(), object.Kind())
				}
			}
		}
	}
}
