package kraph

import (
	"fmt"
	"strings"
	"testing"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/k8s"
	"github.com/milosgajdos/kraph/query"
)

func TestNewKraph(t *testing.T) {
	client, err := k8s.NewMockClient()
	if err != nil {
		t.Fatalf("failed to create API client: %s", err)
	}

	k, err := New(client)
	if err != nil {
		t.Fatalf("failed to create kraph: %v", err)
	}

	obj := k8s.NewMockObject("foo", "bar", "fobar", "randomid")

	node := k.NewNode(obj)
	if node == nil {
		t.Fatal("failed to add new node")
	}

	expCount := 1
	if nodeCount := k.Nodes().Len(); nodeCount != expCount {
		t.Errorf("expected %d nodes,: got:%d", expCount, nodeCount)
	}

	obj = k8s.NewMockObject("foo2", "bar2", "fobar", "randomid2")

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
	client, err := k8s.NewMockClient()
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

func TestQueryNode(t *testing.T) {
	client, err := k8s.NewMockClient()
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

	nodes, err := k.QueryNode()
	if err != nil {
		t.Errorf("failed to query all nodes: %v", err)
	}

	if len(nodes) != g.Nodes().Len() {
		t.Errorf("invalid number of nodes returned. Expected: %d, got: %d", g.Nodes().Len(), len(nodes))
	}

	oddKindNodes, err := k.QueryNode(query.Kind(k8s.MockOddKind))
	if err != nil {
		t.Errorf("failed to query node of kind %s: %v", k8s.MockOddKind, err)
	}

	for _, node := range oddKindNodes {
		obj := node.metadata["object"].(api.Object)
		if !strings.EqualFold(obj.Kind(), k8s.MockOddKind) {
			t.Errorf("expected kind: %s, got %s", k8s.MockOddKind, obj.Kind())
		}
	}

	oddNsOddKindNodes, err := k.QueryNode(query.Kind(k8s.MockOddKind), query.Namespace(k8s.MockOddNs))
	if err != nil {
		t.Errorf("failed to query ns/kind %s/%s: %v", k8s.MockOddNs, k8s.MockOddKind, err)
	}

	for _, node := range oddNsOddKindNodes {
		obj := node.metadata["object"].(api.Object)
		if !strings.EqualFold(obj.Kind(), k8s.MockOddKind) || !strings.EqualFold(obj.Namespace(), k8s.MockOddNs) {
			t.Errorf("expected ns/kind %s/%s, got %s/%s", k8s.MockOddNs, k8s.MockOddKind, obj.Namespace(), obj.Kind())
		}
	}

	name := fmt.Sprintf("%s-%d", k8s.MockAPIOddRes, 1)
	oddNode, err := k.QueryNode(query.Kind(k8s.MockOddKind), query.Namespace(k8s.MockOddNs), query.Name(name))
	if err != nil {
		t.Errorf("failed to query ns/kind/node %s/%s/%s: %v", k8s.MockOddKind, k8s.MockOddNs, name, err)
	}

	expCount := 1
	if len(oddNode) != expCount {
		t.Errorf("expected to find %d node, got: %d", expCount, len(oddNode))
	}
}

func TestQueryEdge(t *testing.T) {
	client, err := k8s.NewMockClient()
	if err != nil {
		t.Fatalf("failed to create API client: %s", err)
	}

	k, err := New(client)
	if err != nil {
		t.Fatalf("failed to create kraph: %v", err)
	}

	_, err = k.Build()
	if err != nil {
		t.Fatalf("failed to build graph: %v", err)
	}
}
