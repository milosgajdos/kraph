package kraph

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/mock"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/memory"
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
	client, err := mock.NewClient()
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
			mock.NewObject("", "pod", "", "", nil),
			[]Filter{func(object api.Object) bool { return object.Kind() == "pod" }},
			false,
		},
		{
			mock.NewObject("", "deployment", "", "", nil),
			[]Filter{func(object api.Object) bool { return object.Kind() == "pod" }},
			true,
		},
		{
			mock.NewObject("name", "", "", "", nil),
			[]Filter{func(object api.Object) bool { return object.Name() == "name" }},
			false,
		},
		{
			mock.NewObject("", "", "", "", nil),
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
