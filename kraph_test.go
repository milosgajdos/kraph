package kraph

import (
	"fmt"
	"strings"
	"testing"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/mock"
	"github.com/milosgajdos/kraph/query"
)

func buildTestKraph() (*Kraph, error) {
	client, err := mock.NewClient()
	if err != nil {
		return nil, err
	}

	k, err := New(client)
	if err != nil {
		return nil, err
	}

	g, err := k.Build()
	if err != nil {
		return nil, err
	}

	if g == nil {
		return nil, fmt.Errorf("nil graph returned")
	}

	return k, nil
}

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
		t.Errorf("failed to build kraph client: %v", err)
	}

	k, err := New(client)
	if err != nil {
		t.Errorf("failed to create new kraph: %v", err)
	}

	g, err := k.Build()
	if err != nil {
		t.Errorf("failed to build new kraph: %v", err)
	}

	if g == nil {
		t.Errorf("nil graph returned")
	}
}

func TestQueryAllNodes(t *testing.T) {
	k, err := buildTestKraph()
	if err != nil {
		t.Fatalf("failed to create new kraph: %v", err)
	}

	nodes, err := k.QueryNode()
	if err != nil {
		t.Errorf("failed to query all nodes: %v", err)
	}

	if len(nodes) != k.Nodes().Len() {
		t.Errorf("invalid number of nodes returned. Expected: %d, got: %d", k.Nodes().Len(), len(nodes))
	}

	for _, nsKinds := range mock.ObjectData {
		for nsKind, names := range nsKinds {
			nsplit := strings.Split(nsKind, "/")
			ns, kind := nsplit[0], nsplit[1]
			for _, name := range names {
				uid := strings.Join([]string{ns, kind, name}, "/")
				nodes, err := k.QueryNode(query.UID(uid))
				if err != nil {
					t.Errorf("error getting node: %s: %v", uid, err)
					continue
				}

				if len(nodes) != 1 {
					t.Errorf("expected single node, got: %d", len(nodes))
					continue
				}

				node := nodes[0]
				object := node.metadata["object"].(api.Object)

				if object.UID().String() != uid {
					t.Errorf("expected node %s, got: %s", uid, object.UID())
				}
			}
		}
	}
}

func TestQueryNodes(t *testing.T) {
	k, err := buildTestKraph()
	if err != nil {
		t.Fatalf("failed to create new kraph: %v", err)
	}

	for _, nsKinds := range mock.ObjectData {
		for nsKind, names := range nsKinds {
			nsplit := strings.Split(nsKind, "/")
			ns, kind := nsplit[0], nsplit[1]

			nodes, err := k.QueryNode(query.Namespace(ns), query.Kind(kind))
			if err != nil {
				t.Errorf("error getting node: %s/%s: %v", ns, kind, err)
				continue
			}

			for _, node := range nodes {
				object := node.metadata["object"].(api.Object)
				if object.Namespace() != ns || object.Kind() != kind {
					t.Errorf("expected: %s/%s, got: %s/%s", ns, kind, object.Namespace(), object.Kind())
				}
			}

			for _, name := range names {
				nodes, err := k.QueryNode(query.Namespace(ns), query.Kind(kind), query.Name(name))
				if err != nil {
					t.Errorf("error getting node: %s/%s/%s: %v", ns, kind, name, err)
					continue
				}

				for _, node := range nodes {
					object := node.metadata["object"].(api.Object)
					if object.Namespace() != ns || object.Kind() != kind {
						t.Errorf("expected: %s/%s/%s, got: %s/%s/%s", ns, kind, name,
							object.Namespace(), object.Kind(), object.Name())
					}
				}
			}
		}
	}
}

func TestQueryAllEdges(t *testing.T) {
	k, err := buildTestKraph()
	if err != nil {
		t.Fatalf("failed to create new kraph: %v", err)
	}

	edges, err := k.QueryEdge()
	if err != nil {
		t.Errorf("failed to query all edges: %v", err)
	}

	expEdges := 0

	for _, rels := range mock.ObjectLinks {
		expEdges += len(rels)
	}

	if len(edges) != expEdges {
		t.Errorf("invalid number of edges returned. Expected: %d, got: %d", expEdges, len(edges))
	}
}

func TestQueryAttrEdges(t *testing.T) {
	k, err := buildTestKraph()
	if err != nil {
		t.Fatalf("failed to create new kraph: %v", err)
	}

	attrs := make(Attrs)

	for _, links := range mock.ObjectLinks {
		for _, relation := range links {
			attrs["relation"] = relation
			edges, err := k.QueryEdge(query.Attrs(attrs))
			if err != nil {
				t.Errorf("failed to query edges with attributes %v: %v", attrs, err)
			}

			for _, edge := range edges {
				if relAttr := edge.GetAttribute("relation"); relAttr != attrs["relation"] {
					t.Errorf("expected relation attribute: %v, got: %v", attrs["relation"], relAttr)
				}
			}
		}
	}
}
