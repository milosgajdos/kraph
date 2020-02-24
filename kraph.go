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
	ErrNotImplemented  = errors.New("not implemented")
	ErrUnknownObject   = errors.New("unknown object")
	ErrInterruptedCall = errors.New("call interrupted")
)

// Node is graph node
type Node struct {
	graph.Node
	Attrs
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
	Attrs
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
	return &Edge{
		WeightedEdge: k.WeightedUndirectedGraph.NewWeightedEdge(from, to, weight),
	}
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
	//fmt.Printf("Linking %s: %d to %s: %d, weight: %f\n", to.Name, to.ID(), from.Name, from.ID(), weight)
	e := k.NewEdge(from, to, weight)
	e.SetAttribute(encoding.Attribute{
		Key:   "weight",
		Value: fmt.Sprintf("%f", weight),
	})
	k.SetWeightedEdge(e)
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

		from := k.NewNode(nodeName(owner.Kind, owner.Name))
		from.SetAttribute(encoding.Attribute{
			Key:   "kind",
			Value: owner.Kind,
		})
		k.AddNode(from)
		k.nodeMap[owner.Kind][owner.UID] = from
		k.addEdge(from, to, weight)
	}
}

// addResNode adds the node item to the kraph graph and links it to its related nodes
func (k *Kraph) addResNode(node *Node, res unstructured.Unstructured) {
	// if the node is NOT in the graph yet, add it in
	if kn := k.Node(node.ID()); kn == nil {
		k.AddNode(node)
	}

	// if the resource is namespaced link it to its namespace
	if ns := res.GetNamespace(); ns != "" {
		node.SetAttribute(encoding.Attribute{
			Key:   "namespace",
			Value: ns,
		})
	}

	k.linkNode(node, res.GetOwnerReferences(), 0.0)
}

// linkResource links API resource to all of its owners
func (k *Kraph) linkResource(res unstructured.Unstructured) {
	kind := strings.ToLower(res.GetKind())
	//fmt.Println("Item kind:", kind, " Item name:", res.GetName(), " Item UID:", res.GetUID(), " Owners:", res.GetOwnerReferences())
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
		k.addResNode(node, res)
		return
	}

	node := k.NewNode(nodeName(kind, res.GetName()))
	node.SetAttribute(encoding.Attribute{
		Key:   "kind",
		Value: kind,
	})
	k.addResNode(node, res)
	k.nodeMap[kind][uid] = node
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
			k.linkResource(item)
		}
	}

	errChan <- err
}
