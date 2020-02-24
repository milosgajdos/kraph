package kraph

import (
	"context"
	"errors"
	"fmt"
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
	ErrNotImplemented  = errors.New("not implemented")
	ErrUnknownObject   = errors.New("unknown object")
	ErrInterruptedCall = errors.New("call interrupted")
)

// Node is graph node
type Node struct {
	graph.Node
	Attributes
	// Name names the node
	Name string
}

// DOTID returns the node's DOT ID.
func (n *Node) DOTID() string {
	return n.Name
}

// SetDOTID sets the node's DOT ID.
func (n *Node) SetDOTID(id string) {
	n.Name = id
}

// Edge is graph edge
type Edge struct {
	graph.WeightedEdge
	Attributes
}

// Kraph is a graph of Kubernetes resources
type Kraph struct {
	*simple.WeightedUndirectedGraph
	// nodeMap maps graph nodes to kubernetes IDs per kind
	nodeMap map[string]map[types.UID]*Node
	// disc is kubernetes discovery client
	disc discovery.DiscoveryInterface
	// dyn is kubernetes dynamic client
	dyn dynamic.Interface
	// Global DOT attributes
	GraphAttrs Attributes
	NodeAttrs  Attributes
	EdgeAttrs  Attributes
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
func (k *Kraph) NewNode(name string) *Node {
	// TODO: should add the node to the graph for better UX
	// as if its not added in, it won't have unique ID
	return &Node{
		Node: k.WeightedUndirectedGraph.NewNode(),
		Name: name,
	}
}

// NewEdge adds a new edge from source node to destination node to the graph
// or returns an existing edge if it already exists
// It will panic if the IDs of the from and to are equal
func (k *Kraph) NewEdge(from, to *Node, weight float64) *Edge {
	if e := k.Edge(from.ID(), to.ID()); e != nil {
		ke, ok := e.(*Edge)
		if !ok {
			return &Edge{
				WeightedEdge: e.(*simple.WeightedEdge),
			}
		}

		return ke
	}

	e := &Edge{
		WeightedEdge: k.WeightedUndirectedGraph.NewWeightedEdge(from, to, weight),
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
func (k *Kraph) addEdge(from, to *Node, weight float64) {
	//fmt.Printf("Linking %s to %s\n", to.Name, from.Name)
	e := k.NewEdge(from, to, weight)
	e.SetAttribute(encoding.Attribute{
		Key:   "weight",
		Value: fmt.Sprintf("%f", weight),
	})
}

// linkNode links the node too all of its owners
func (k *Kraph) linkNode(to *Node, owners []metav1.OwnerReference, weight float64) {
	for _, owner := range owners {
		kind := owner.Kind
		if k.nodeMap[kind] == nil {
			k.nodeMap[kind] = make(map[types.UID]*Node)
		}
		if from, ok := k.nodeMap[owner.Kind][owner.UID]; ok {
			k.addEdge(from, to, weight)
			continue
		}

		from := k.NewNode(owner.Name)
		from.SetAttribute(encoding.Attribute{
			Key:   "kind",
			Value: owner.Kind,
		})
		k.AddNode(from)
		k.addEdge(from, to, 0.0)
	}
}

// addNodeItem adds the node item to the kraph graph and links it to its related nodes
func (k *Kraph) addNodeItem(node *Node, item unstructured.Unstructured) {
	if kn := k.Node(node.ID()); kn == nil {
		k.AddNode(node)
	}

	// if the item is namespaced link it to its namespace
	if ns := item.GetNamespace(); ns != "" {
		kind := item.GetKind()
		if k.nodeMap[kind] == nil {
			k.nodeMap[kind] = make(map[types.UID]*Node)
		}
		//fmt.Println("Item", item.GetName(), "is namespaced to", ns)
		nsNode, ok := k.nodeMap[kind][types.UID(ns)]
		if !ok {
			nsNode = k.NewNode(ns)
			node.SetAttribute(encoding.Attribute{
				Key:   "kind",
				Value: "namespace",
			})
			k.AddNode(nsNode)
		}
		k.addEdge(node, nsNode, 0.0)
	}

	k.linkNode(node, item.GetOwnerReferences(), 0.0)
}

// linkItem links an API resource item to its owners
func (k *Kraph) linkItem(item unstructured.Unstructured) {
	kind := item.GetKind()
	if k.nodeMap[kind] == nil {
		k.nodeMap[kind] = make(map[types.UID]*Node)
	}
	if node, ok := k.nodeMap[kind][item.GetUID()]; ok {
		k.addNodeItem(node, item)
		return
	}

	node := k.NewNode(item.GetName())
	node.SetAttribute(encoding.Attribute{
		Key:   "kind",
		Value: item.GetKind(),
	})
	k.addNodeItem(node, item)
	k.nodeMap[kind][item.GetUID()] = node
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

		//fmt.Println("Discovered", result.api, "objects:", len(result.items))

		for _, item := range result.items {
			k.linkItem(item)
		}
	}

	errChan <- err
}
