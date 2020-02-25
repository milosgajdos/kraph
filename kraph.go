package kraph

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

var (
	// ErrNotImplemented is returned when some functionality has not been implemented
	ErrNotImplemented = errors.New("not implemented")
	// ErrUnknownObject is returned when an unknown object has been been provided
	ErrUnknownObject = errors.New("unknown object")
)

// Kraph is a graph of Kubernetes resources
type Kraph struct {
	*simple.WeightedUndirectedGraph
	// nodeMap maps API resources as graph nodes key-ed over Kind
	nodeMap map[string]map[types.UID]*Node
	// disc is kubernetes discovery client
	disc discovery.DiscoveryInterface
	// dyn is kubernetes dynamic client
	dyn dynamic.Interface
	// Global DOT attributes
	GraphAttrs Attrs
	NodeAttrs  Attrs
	EdgeAttrs  Attrs
}

// New creates new Kraph and returns it
//func New(client kubernetes.Interface) (*Kraph, error) {
func New(discover discovery.DiscoveryInterface, dynamic dynamic.Interface) (*Kraph, error) {
	return &Kraph{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(0.0, 0.0),
		nodeMap:                 make(map[string]map[types.UID]*Node),
		disc:                    discover,
		dyn:                     dynamic,
	}, nil
}

// NewNode creates new kraph node and returns it
func (k *Kraph) NewNode(name string, attrs ...encoding.Attribute) graph.Node {
	// as if its not added in, it won't have unique ID
	n := &Node{
		Node: k.WeightedUndirectedGraph.NewNode(),
		Name: name,
	}

	for _, attr := range attrs {
		n.SetAttribute(attr)
	}

	return n
}

// NewEdge adds a new edge from source node to destination node to the graph
// or returns an existing edge if it already exists
// It will panic if the IDs of the from and to are equal
func (k *Kraph) NewEdge(from, to graph.Node, weight float64, attrs ...encoding.Attribute) graph.Edge {
	if e := k.Edge(from.ID(), to.ID()); e != nil {
		return e
	}

	e := &Edge{
		from:   from.(*Node),
		to:     to.(*Node),
		weight: weight,
	}

	for _, attr := range attrs {
		e.SetAttribute(attr)
	}

	k.SetWeightedEdge(e)

	return e
}

// DOTID returns the graph's DOT ID.
func (k *Kraph) DOTID() string {
	return "kraph"
}

// DOTAttributers returns the global DOT kraph attributers
func (k *Kraph) DOTAttributers() (graph, node, edge encoding.Attributer) {
	return k.GraphAttrs, k.NodeAttrs, k.EdgeAttrs
}

// DOT returns the GrapViz dot representation of kraph
func (k *Kraph) DOT() (string, error) {
	b, err := dot.Marshal(k, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode kraph into dot: %v", err)
	}

	return string(b), nil
}

// Build builds the kubernetes resource graph for a given namespace
// If the namespace is empty string all namespaces are assumed
func (k *Kraph) Build(ctx context.Context, namespace string) error {
	api, err := k.discoverAPI(ctx)
	if err != nil {
		return fmt.Errorf("failed discovering kubernetes API: %w", err)
	}

	return k.buildGraph(ctx, api, namespace)
}

// GetNodesWithAttr returns a slice of nodes with the given attribute set
// If it does not find any matching nodes it returns empty slice.
// It returns error if attribute key is empty string,
func (k *Kraph) GetNodesWithAttr(attr encoding.Attribute) ([]*Node, error) {
	var nodes []*Node

	if attr.Key == "" {
		return nil, ErrAttrKeyInvalid
	}

	found := false
	for _, node := range graph.NodesOf(k.Nodes()) {
		n := node.(*Node)
		if val := n.Get(attr.Key); val != "" {
			// attribute key exists; check its value
			switch attr.Value {
			case "*":
				found = true
			case val:
				found = true
			default:
				// continue
			}
		}
		if found {
			nodes = append(nodes, n)
			found = false
		}
	}

	return nodes, nil
}

// GetEdgesWithAttr returns a slice of Edges with the given attribute set
// If it does not find any matching edges it returns empty slice.
// It returns error if attribute key is empty string,
func (k *Kraph) GetEdgesWithAttr(attr encoding.Attribute) ([]*Edge, error) {
	var edges []*Edge

	if attr.Key == "" {
		return nil, ErrAttrKeyInvalid
	}

	found := false
	for _, edge := range graph.EdgesOf(k.Edges()) {
		e := edge.(*Edge)
		if val := e.Get(attr.Key); val != "" {
			// attribute key exists; check its value
			switch attr.Value {
			case "*":
				found = true
			case val:
				found = true
			default:
				// continue
			}
		}
		if found {
			edges = append(edges, e)
			found = false
		}
	}

	return edges, nil
}

// discoverAPI discovers all available kubernetes API resource groups and returns them
// It returns error if it fails to retrieve the resources of if it fails to parse their versions
func (k *Kraph) discoverAPI(ctx context.Context) (*API, error) {
	srvPrefResList, err := k.disc.ServerPreferredResources()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch api groups from kubernetes: %w", err)
	}

	api := &API{
		resourceMap: make(map[string][]Resource),
	}

	for _, srvPrefRes := range srvPrefResList {
		gv, err := schema.ParseGroupVersion(srvPrefRes.GroupVersion)
		if err != nil {
			return nil, fmt.Errorf("failed parsing %s into GroupVersion: %w", srvPrefRes.GroupVersion, err)
		}

		for _, apiResource := range srvPrefRes.APIResources {
			if !provides(apiResource.Verbs, "list") {
				continue
			}

			resource := Resource{
				ar: apiResource,
				gv: gv,
			}

			api.resources = append(api.resources, resource)
			for _, path := range resource.Paths() {
				api.resourceMap[path] = append(api.resourceMap[path], resource)
			}
		}
	}

	return api, nil
}

type result struct {
	api   string
	items []unstructured.Unstructured
	err   error
}

// buildGraph builds a graph of api resources in a given namespace
// if the namespace is empty it queries API groups across all namespaces
// It returns error if any of the API calls fails with error.
func (k *Kraph) buildGraph(ctx context.Context, api *API, ns string) error {
	// TODO: we should take into account the context when firing goroutines
	var wg sync.WaitGroup

	resChan := make(chan result, 250)
	doneChan := make(chan struct{})

	for _, resource := range api.Resources() {
		// if all namespaces are scanned and the API resource is namespaced, skip
		if ns != "" && !resource.ar.Namespaced {
			continue
		}

		gvResClient := k.dyn.Resource(schema.GroupVersionResource{
			Group:    resource.gv.Group,
			Version:  resource.gv.Version,
			Resource: resource.ar.Name,
		})

		var client dynamic.ResourceInterface
		switch ns {
		case "":
			client = gvResClient
		default:
			client = gvResClient.Namespace(ns)
		}

		wg.Add(1)
		go func(r Resource) {
			defer wg.Done()
			var cont string
			for {
				res, err := client.List(metav1.ListOptions{
					Limit:    100,
					Continue: cont,
				})
				select {
				case resChan <- result{api: r.ar.Name, items: res.Items, err: err}:
				case <-doneChan:
					return
				}
				cont = res.GetContinue()
				if cont == "" {
					break
				}
			}
		}(resource)
	}

	errChan := make(chan error, 1)
	go k.processResults(resChan, doneChan, errChan)

	wg.Wait()
	close(resChan)

	err := <-errChan

	return err
}

// addEdge creates an between from and to nodes with given weight
func (k *Kraph) addEdge(from, to graph.Node, weight float64) {
	attrs := []encoding.Attribute{
		encoding.Attribute{Key: "weight", Value: fmt.Sprintf("%f", weight)},
	}

	k.NewEdge(from, to, weight, attrs...)
}

// linkNode links the node too all of its owners
func (k *Kraph) linkNode(from graph.Node, owners []metav1.OwnerReference, weight float64) {
	for _, owner := range owners {
		kind := owner.Kind
		if k.nodeMap[kind] == nil {
			k.nodeMap[kind] = make(map[types.UID]*Node)
		}
		if to, ok := k.nodeMap[owner.Kind][owner.UID]; ok {
			if e := k.Edge(from.ID(), to.ID()); e == nil {
				k.addEdge(from, to, weight)
			}
			continue
		}

		name := nodeName(owner.Kind, owner.Name)
		attrs := []encoding.Attribute{
			encoding.Attribute{Key: "kind", Value: owner.Kind},
		}
		to := k.NewNode(name, attrs...)
		k.AddNode(to)
		k.nodeMap[owner.Kind][owner.UID] = to.(*Node)
		k.addEdge(from, to, weight)
	}
}

// linkResNode adds the node item to the kraph graph and links it to its related nodes
func (k *Kraph) linkResNode(node graph.Node, res unstructured.Unstructured) {
	// if the node is NOT in the graph yet, add it in
	if kn := k.Node(node.ID()); kn == nil {
		k.AddNode(node)
	}

	// if the resource is namespaced link it to its namespace
	if ns := res.GetNamespace(); ns != "" {
		knode := node.(*Node)
		knode.SetAttribute(encoding.Attribute{
			Key:   "namespace",
			Value: ns,
		})
	}

	k.linkNode(node, res.GetOwnerReferences(), 0.0)
}

// linkResource links API resource to all of its owners
func (k *Kraph) linkResource(res unstructured.Unstructured) {
	kind := strings.ToLower(res.GetKind())
	if k.nodeMap[kind] == nil {
		k.nodeMap[kind] = make(map[types.UID]*Node)
	}

	// TODO: Some API resrouces like ComponentStatus do NOT have UID set o_O
	// Let's concatenate the name of the resource with the kind for now
	uid := res.GetUID()
	if uid == "" {
		uid = types.UID(kind + res.GetName())
	}
	if node, ok := k.nodeMap[kind][uid]; ok {
		k.linkResNode(node, res)
		return
	}

	name := nodeName(kind, res.GetName())
	attrs := []encoding.Attribute{
		encoding.Attribute{Key: "kind", Value: kind},
	}
	node := k.NewNode(name, attrs...)
	k.AddNode(node)
	k.linkResNode(node, res)
	k.nodeMap[kind][uid] = node.(*Node)
}

// processResults process API calls request results
// Ot builds the undirected weighted graph from the received results
func (k *Kraph) processResults(resChan <-chan result, doneChan chan struct{}, errChan chan<- error) {
	var err error
	for result := range resChan {
		if result.err != nil {
			err = result.err
			close(doneChan)
			break
		}

		for _, item := range result.items {
			k.linkResource(item)
		}
	}

	errChan <- err
}
