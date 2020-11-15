package kraph

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gen"
	"github.com/milosgajdos/kraph/pkg/store"
	"github.com/milosgajdos/kraph/pkg/store/memory"
)

const (
	resPath = "pkg/api/gen/seeds/resources.yaml"
	objPath = "pkg/api/gen/seeds/objects.yaml"
)

func TestNewKraph(t *testing.T) {
	s, err := memory.NewStore("default", store.Options{})
	if err != nil {
		t.Fatalf("failed to create memory store: %v", err)
	}

	k, err := New(s)
	if err != nil {
		t.Fatalf("failed to create kraph: %v", err)
	}

	if k == nil {
		t.Fatal("got nil kraph")
	}
}

func TestBuild(t *testing.T) {
	client, err := gen.NewMockClient(resPath, objPath)
	if err != nil {
		t.Errorf("failed to build mock client: %v", err)
	}

	m, err := memory.NewStore("memory", store.Options{})
	if err != nil {
		t.Fatalf("failed to create memory store: %v", err)
	}

	k, err := New(m)
	if err != nil {
		t.Fatalf("failed to create kraph: %v", err)
	}

	if err := k.Build(client); err != nil {
		t.Errorf("failed to build graph: %v", err)
	}
}

func TestStore(t *testing.T) {
	m, err := memory.NewStore("memory", store.Options{})
	if err != nil {
		t.Fatalf("failed to create memory store: %v", err)
	}

	k, err := New(m)
	if err != nil {
		t.Errorf("failed to build mock client: %v", err)
	}

	s := k.Store()

	if !reflect.DeepEqual(s, m) {
		t.Errorf("expected store: %#v, got: %#v", m, s)
	}
}

func TestSkipGraph(t *testing.T) {
	tests := []struct {
		object   api.Object
		filters  []Filter
		expected bool
	}{
		{
			gen.NewMockObject("", "", "", gen.NewMockResource("", "pod", "", "", false, api.Options{}), api.Options{}),
			[]Filter{func(object api.Object) bool { return object.Resource().Kind() == "pod" }},
			false,
		},
		{
			gen.NewMockObject("", "", "", gen.NewMockResource("", "deployment", "", "", false, api.Options{}), api.Options{}),
			[]Filter{func(object api.Object) bool { return object.Resource().Kind() == "pod" }},
			true,
		},
		{
			gen.NewMockObject("", "name", "", nil, api.Options{}),
			[]Filter{func(object api.Object) bool { return object.Name() == "name" }},
			false,
		},
		{
			gen.NewMockObject("", "", "", nil, api.Options{}),
			[]Filter{},
			false,
		},
	}
	for _, test := range tests {
		if skip(test.object, test.filters...) != test.expected {
			t.Errorf("expected: %v, got: %v, for: %#v", test.expected, !test.expected, test.object)

		}
	}
}
