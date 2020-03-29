package kraph

import (
	"testing"

	"github.com/milosgajdos/kraph/api/k8s"
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

	obj := k8s.NewMockObject("foo", "bar", "fobar")

	node := k.NewNode(obj)
	if node == nil {
		t.Fatal("failed to add new node")
	}

	expCount := 1
	if nodeCount := k.Nodes().Len(); nodeCount != expCount {
		t.Errorf("expected %d nodes,: got:%d", expCount, nodeCount)
	}

	obj = k8s.NewMockObject("foo2", "bar2", "fobar")

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

//func TestNodeAttributes(t *testing.T) {
//	disc := testclient.NewSimpleClientset().Discovery()
//	dyn := testdynclient.NewSimpleDynamicClient(runtime.NewScheme())
//	k, err := New(disc, dyn)
//	if err != nil {
//		t.Fatalf("failed creating new kraph: %v", err)
//	}
//
//	// add 3 foo nodes
//	fooCount := 3
//	addNodes(k, "foo", fooCount)
//
//	// add 2 foo nodes
//	barCount := 2
//	addNodes(k, "bar", barCount)
//
//	nodes, err := k.GetNodesWithAttr(encoding.Attribute{Key: "foo", Value: "*"})
//	if err != nil {
//		t.Errorf("failed adding foo nodes: %v", err)
//	}
//
//	if len(nodes) != fooCount {
//		t.Errorf("invalid number of foo nodes returned. expected: %d, got: %d", fooCount, len(nodes))
//	}
//
//	if _, err := k.GetNodesWithAttr(encoding.Attribute{Key: "", Value: "*"}); err != ErrAttrKeyInvalid {
//		t.Errorf("expected to fail with %v, got: %v", ErrAttrKeyInvalid, err)
//	}
//
//	nodes, err = k.GetNodesWithAttr(encoding.Attribute{Key: "foo", Value: ""})
//	if err != nil {
//		t.Errorf("failed querying node attributes: %v", err)
//	}
//
//	if len(nodes) != 0 {
//		t.Errorf("incorrect number of nodes returned, expected: %d, got: %d", 0, len(nodes))
//	}
//}
//
//func TestEdgeAttributes(t *testing.T) {
//	disc := testclient.NewSimpleClientset().Discovery()
//	dyn := testdynclient.NewSimpleDynamicClient(runtime.NewScheme())
//	k, err := New(disc, dyn)
//	if err != nil {
//		t.Fatalf("failed creating new kraph: %v", err)
//	}
//
//	// add 5 foo nodes
//	fooCount := 5
//	addNodes(k, "foo", fooCount)
//
//	// add bar edges between 1-2 and 2-4
//	attr := encoding.Attribute{Key: "bar", Value: "foo"}
//	nodes := graph.NodesOf(k.Nodes())
//
//	k.NewEdge(nodes[0], nodes[1], 0.0, attr)
//	k.NewEdge(nodes[1], nodes[3], 0.0, attr)
//
//	edges, err := k.GetEdgesWithAttr(encoding.Attribute{Key: "bar", Value: "*"})
//	if err != nil {
//		t.Errorf("failed getting bar edges: %v", err)
//	}
//
//	if len(edges) != 2 {
//		t.Errorf("invalid number of foo nodes returned. expected: %d, got: %d", 2, len(edges))
//	}
//
//	if _, err := k.GetEdgesWithAttr(encoding.Attribute{Key: "", Value: "*"}); err != ErrAttrKeyInvalid {
//		t.Errorf("expected to fail with %v, got: %v", ErrAttrKeyInvalid, err)
//	}
//
//	edges, err = k.GetEdgesWithAttr(encoding.Attribute{Key: "bar", Value: ""})
//	if err != nil {
//		t.Errorf("failed querying edge attributes: %v", err)
//	}
//
//	if len(edges) != 0 {
//		t.Errorf("incorrect number of edges returned, expected: %d, got: %d", 0, len(edges))
//	}
//}
