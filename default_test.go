package kraph

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/kraph/api/mock"
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

	k, err := New()
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
	m := memory.New("memory")

	k, err := New(Store(m))
	if err != nil {
		t.Errorf("failed to build mock client: %v", err)
	}

	s := k.Store()

	if !reflect.DeepEqual(s, m) {
		t.Errorf("expected store: %#v, got: %#v", m, s)
	}
}
