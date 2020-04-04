package kraph

import (
	"testing"

	"github.com/milosgajdos/kraph/api/mock"
)

func TestNewKraph(t *testing.T) {
	client, err := mock.NewClient()
	if err != nil {
		t.Fatalf("failed to create API client: %s", err)
	}

	k, err := New(client)
	if err != nil {
		t.Fatalf("failed to create kraph: %v", err)
	}

	obj := mock.NewObject("foo", "bar", "fobar", "randomid", nil)

	node := k.NewNode(obj)
	if node == nil {
		t.Fatal("failed to add new node")
	}

	expCount := 1
	if nodeCount := k.Nodes().Len(); nodeCount != expCount {
		t.Errorf("expected %d nodes,: got:%d", expCount, nodeCount)
	}

	obj = mock.NewObject("foo2", "bar2", "fobar", "randomid2", nil)

	node2 := k.NewNode(obj)
	if node2 == nil {
		t.Fatal("failed to add new node")
	}

	edge := k.NewEdge(node, node2)
	if edge == nil {
		t.Fatal("failed to create edge")
	}

	expCount = 1
	if edgeCount := k.Edges().Len(); edgeCount != expCount {
		t.Errorf("expected: %d edges, got: %d", expCount, edgeCount)
	}

	g, n, e := k.DOTAttributers()
	if len(g.Attributes()) != 0 || len(n.Attributes()) != 0 || len(e.Attributes()) != 0 {
		t.Errorf("invalid DOT attributes, expected no attributes, got: %d, %d, %d",
			len(g.Attributes()), len(n.Attributes()), len(e.Attributes()))
	}

	dot, err := k.DOT()
	if err != nil {
		t.Errorf("failed getting DOT graph: %v", err)
	}

	if dot == "" {
		t.Errorf("empty DOT graph returned, expected non-empty graph")
	}
}

func TestBuild(t *testing.T) {
	client, err := mock.NewClient()
	if err != nil {
		t.Fatalf("failed to create API client: %s", err)
	}

	k, err := New(client)
	if err != nil {
		t.Fatalf("failed to create kraph: %v", err)
	}

	g, err := k.Build()
	if err != nil {
		t.Fatalf("failed to build graph: %v", err)
	}

	if g == nil {
		t.Fatal("nil graph returned")
	}
}
