package kraph

import (
	"testing"

	"gonum.org/v1/gonum/graph/simple"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestNode(t *testing.T) {
	g := simple.NewWeightedUndirectedGraph(0.0, 0.0)
	name := "foo"
	node := &Node{
		Node: g.NewNode(),
		Name: name,
	}

	if dotID := node.DOTID(); dotID != name {
		t.Errorf("expected: %s, go: %s", name, dotID)
	}

	id := "bar"
	node.SetDOTID(id)

	if dotID := node.DOTID(); dotID != id {
		t.Errorf("expected: %s, go: %s", id, dotID)
	}
}

func TestNewKraph(t *testing.T) {
	k, err := New(testclient.NewSimpleClientset())
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

	existingEdge := k.NewEdge(node, node2, 0.0)
	from, to := existingEdge.From().ID(), existingEdge.To().ID()
	if from != edge.From().ID() || to != edge.To().ID() {
		t.Errorf("invalid edge returned. expected [from, to]: [%d, %d], got: [%d, %d]",
			edge.From().ID(), edge.To().ID(), from, to)
	}

	// add a new pair of nodes
	fromNode := k.NewNode("foo2")
	k.AddNode(fromNode)
	toNode := k.NewNode("bar2")
	k.AddNode(toNode)

	// add graph.Edge instead of kraph.Edge
	graphEdge := k.WeightedUndirectedGraph.NewWeightedEdge(fromNode, toNode, 0.1)
	k.SetWeightedEdge(graphEdge)
	// we should get back kraph.Edge instead of graph.Edge
	existingEdge = k.NewEdge(fromNode, toNode, 0.0)
	from, to = graphEdge.From().ID(), graphEdge.To().ID()
	if from != graphEdge.From().ID() || to != graphEdge.To().ID() {
		t.Errorf("invalid edge returned. expected [from, to]: [%d, %d], got: [%d, %d]",
			graphEdge.From().ID(), graphEdge.To().ID(), from, to)
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
	k, err := New(testclient.NewSimpleClientset())
	if err != nil {
		t.Fatalf("failed creating new kraph: %v", err)
	}

	if err := k.Build(); err != nil {
		t.Errorf("failed to build kraph: %v", err)
	}
}
