package kraph

import (
	"context"
	"fmt"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"k8s.io/apimachinery/pkg/runtime"
	testdynclient "k8s.io/client-go/dynamic/fake"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func newKraph() (*Kraph, error) {
	disc := testclient.NewSimpleClientset().Discovery()
	dyn := testdynclient.NewSimpleDynamicClient(runtime.NewScheme())
	return New(disc, dyn)
}

func addNodes(k *Kraph, attr string, count int) {
	nodes := make([]*Node, count)

	for i := 0; i < count; i++ {
		name := fmt.Sprintf("%d", i)
		attr := encoding.Attribute{Key: attr, Value: name}
		node := k.NewNode(name, attr)
		k.AddNode(node)
		nodes = append(nodes, node.(*Node))
	}
}

func TestNewKraph(t *testing.T) {
	k, err := newKraph()
	if err != nil {
		t.Fatalf("failed creating new kraph: %v", err)
	}

	node := k.NewNode("foo")
	if node == nil {
		t.Errorf("failed to create new kraph node")
	}

	if nodeCount := k.Nodes().Len(); nodeCount != 0 {
		t.Errorf("invalid kraph nodes, expected: %d, got:%d", 0, nodeCount)
	}

	k.AddNode(node)
	if nodeCount := k.Nodes().Len(); nodeCount != 1 {
		t.Errorf("invalid number of kraph nodes, expected: %d, got:%d", 1, nodeCount)
	}

	node2 := k.NewNode("bar")
	k.AddNode(node2)

	edge := k.NewEdge(node, node2, 0.0)
	if edge == nil {
		t.Errorf("failed to create new kraph edge")
	}

	if edgeCount := k.Edges().Len(); edgeCount != 1 {
		t.Errorf("invalid number of kraph edges, expected: %d, got:%d", 1, edgeCount)
	}

	g, n, e := k.DOTAttributers()
	if len(g.Attributes()) != 0 || len(n.Attributes()) != 0 || len(e.Attributes()) != 0 {
		t.Errorf("invalid DOT attributes, expected 0 attributes, got: %d, %d, %d",
			len(g.Attributes()), len(n.Attributes()), len(e.Attributes()))
	}

	dotKraph, err := k.DOT()
	if err != nil {
		t.Errorf("failed getting DOT graph: %v", err)
	}

	if dotKraph == "" {
		t.Errorf("empty DOT graph returned, expected non-empty graph")
	}
}

func TestBuild(t *testing.T) {
	k, err := newKraph()
	if err != nil {
		t.Fatalf("failed creating new kraph: %v", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := k.Build(ctx, ""); err != nil {
		t.Errorf("failed to build kraph: %v", err)
	}
}

func TestNodeAttributes(t *testing.T) {
	disc := testclient.NewSimpleClientset().Discovery()
	dyn := testdynclient.NewSimpleDynamicClient(runtime.NewScheme())
	k, err := New(disc, dyn)
	if err != nil {
		t.Fatalf("failed creating new kraph: %v", err)
	}

	// add 3 foo nodes
	fooCount := 3
	addNodes(k, "foo", fooCount)

	// add 2 foo nodes
	barCount := 2
	addNodes(k, "bar", barCount)

	nodes, err := k.GetNodesWithAttr(encoding.Attribute{Key: "foo", Value: "*"})
	if err != nil {
		t.Errorf("failed adding foo nodes: %v", err)
	}

	if len(nodes) != fooCount {
		t.Errorf("invalid number of foo nodes returned. expected: %d, got: %d", fooCount, len(nodes))
	}

	if _, err := k.GetNodesWithAttr(encoding.Attribute{Key: "", Value: "*"}); err != ErrAttrKeyInvalid {
		t.Errorf("expected to fail with %v, got: %v", ErrAttrKeyInvalid, err)
	}

	nodes, err = k.GetNodesWithAttr(encoding.Attribute{Key: "foo", Value: ""})
	if err != nil {
		t.Errorf("failed querying node attributes: %v", err)
	}

	if len(nodes) != 0 {
		t.Errorf("incorrect number of nodes returned, expected: %d, got: %d", 0, len(nodes))
	}
}

func TestEdgeAttributes(t *testing.T) {
	disc := testclient.NewSimpleClientset().Discovery()
	dyn := testdynclient.NewSimpleDynamicClient(runtime.NewScheme())
	k, err := New(disc, dyn)
	if err != nil {
		t.Fatalf("failed creating new kraph: %v", err)
	}

	// add 5 foo nodes
	fooCount := 5
	addNodes(k, "foo", fooCount)

	// add bar edges between 1-2 and 2-4
	attr := encoding.Attribute{Key: "bar", Value: "foo"}
	nodes := graph.NodesOf(k.Nodes())

	k.NewEdge(nodes[0], nodes[1], 0.0, attr)
	k.NewEdge(nodes[1], nodes[3], 0.0, attr)

	edges, err := k.GetEdgesWithAttr(encoding.Attribute{Key: "bar", Value: "*"})
	if err != nil {
		t.Errorf("failed getting bar edges: %v", err)
	}

	if len(edges) != 2 {
		t.Errorf("invalid number of foo nodes returned. expected: %d, got: %d", 2, len(edges))
	}

	if _, err := k.GetEdgesWithAttr(encoding.Attribute{Key: "", Value: "*"}); err != ErrAttrKeyInvalid {
		t.Errorf("expected to fail with %v, got: %v", ErrAttrKeyInvalid, err)
	}

	edges, err = k.GetEdgesWithAttr(encoding.Attribute{Key: "bar", Value: ""})
	if err != nil {
		t.Errorf("failed querying edge attributes: %v", err)
	}

	if len(edges) != 0 {
		t.Errorf("incorrect number of edges returned, expected: %d, got: %d", 0, len(edges))
	}
}
