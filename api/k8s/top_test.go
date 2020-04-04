package k8s

import (
	"strings"
	"testing"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/mock"
	"github.com/milosgajdos/kraph/query"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func newTestTop() *Top {
	top := newTop()

	for resName, meta := range mock.Resources {
		groups := mock.ResourceData[resName]["groups"]
		versions := mock.ResourceData[resName]["versions"]
		for _, group := range groups {
			for _, version := range versions {
				gv := strings.Join([]string{group, version}, "/")
				if gvObject, ok := mock.ObjectData[gv]; ok {
					ns := meta["ns"]
					if len(ns) == 0 {
						ns = api.NamespaceNan
					}

					nsKind := strings.Join([]string{ns, meta["kind"]}, "/")
					if names, ok := gvObject[nsKind]; ok {
						for _, name := range names {
							uid := strings.Join([]string{ns, meta["kind"], name}, "/")
							object := &Object{
								name:  name,
								kind:  meta["kind"],
								ns:    ns,
								uid:   &UID{uid: uid},
								raw:   unstructured.Unstructured{},
								links: make(map[string]*Relation),
							}
							top.Add(object)
						}
					}
				}
			}
		}
	}

	return top
}

func TestObjects(t *testing.T) {
	top := newTestTop()

	objects := top.Objects()
	if len(objects) == 0 {
		t.Errorf("no objects found")
	}
}

func TestGetUID(t *testing.T) {
	top := newTestTop()

	for _, nsKinds := range mock.ObjectData {
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
	top := newTestTop()

	objects, err := top.Get()
	if err != nil {
		t.Errorf("error getting all namespace objects: %v", err)
	}

	if len(objects) == 0 {
		t.Errorf("no mamespaced objects returned")
	}

	for _, nsKinds := range mock.ObjectData {
		for nsKind, names := range nsKinds {
			nsplit := strings.Split(nsKind, "/")
			ns, kind := nsplit[0], nsplit[1]

			objects, err := top.Get(query.Namespace(ns))
			if err != nil {
				t.Errorf("error getting namespace %s objects: %v", ns, err)
				continue
			}

			for _, object := range objects {
				if object.Namespace() != ns {
					t.Errorf("expected: namespace %s, got: %s", ns, object.Namespace())
				}
			}

			objects, err = top.Get(query.Namespace(ns), query.Kind(kind))
			if err != nil {
				t.Errorf("error getting objects: %s/%s: %v", ns, kind, err)
				continue
			}

			for _, object := range objects {
				if object.Namespace() != ns || object.Kind() != kind {
					t.Errorf("expected: %s/%s, got: %s/%s", ns, kind, object.Namespace(), object.Kind())
				}
			}

			for _, name := range names {
				objects, err = top.Get(query.Namespace(ns), query.Kind(kind), query.Name(name))
				if err != nil {
					t.Errorf("error getting objects: %s/%s/%s: %v", ns, kind, name, err)
					continue
				}

				for _, object := range objects {
					if object.Namespace() != ns || object.Kind() != kind || object.Name() != name {
						t.Errorf("expected: %s/%s/%s, got: %s/%s/%s", ns, kind, name,
							object.Namespace(), object.Kind(), object.Name())
					}
				}
			}
		}
	}
}
