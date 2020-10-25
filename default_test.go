package kraph

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/gen"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/memory"
)

const (
	resPath = "seeds/resources.yaml"
	objPath = "seeds/objects.yaml"
)

func TestNewKraph(t *testing.T) {
	k, err := New()
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

	k, err := New(Store(m))
	if err != nil {
		t.Fatalf("failed to create kraph: %v", err)
	}

	g, err := k.Build(client)
	if err != nil {
		t.Errorf("failed to build graph: %v", err)
	}

	if g == nil {
		t.Errorf("nil graph returned")
	}
}

func TestStore(t *testing.T) {
	m, err := memory.NewStore("memory", store.Options{})
	if err != nil {
		t.Fatalf("failed to create memory store: %v", err)
	}

	k, err := New(Store(m))
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
			gen.NewMockObject("", "", "", gen.NewMockResource("", "pod", "", "", false)),
			[]Filter{func(object api.Object) bool { return object.Resource().Kind() == "pod" }},
			false,
		},
		{
			gen.NewMockObject("", "", "", gen.NewMockResource("", "deployment", "", "", false)),
			[]Filter{func(object api.Object) bool { return object.Resource().Kind() == "pod" }},
			true,
		},
		{
			gen.NewMockObject("", "name", "", nil),
			[]Filter{func(object api.Object) bool { return object.Name() == "name" }},
			false,
		},
		{
			gen.NewMockObject("", "", "", nil),
			[]Filter{},
			false,
		},
	}
	for _, test := range tests {
		if skipGraph(test.object, test.filters...) != test.expected {
			t.Errorf("expected: %v, got: %v, for: %#v", test.expected, !test.expected, test.object)

		}
	}
}
