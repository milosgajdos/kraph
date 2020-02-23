package kraph

import (
	"errors"
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrUnknownObject  = errors.New("unknown object")
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
	// nodeMap maps graph nodes to kubernetes IDs
	nodeMap map[types.UID]*Node
	// client is kubernetes clientset
	client kubernetes.Interface
	// Global DOT attributes
	GraphAttrs Attributes
	NodeAttrs  Attributes
	EdgeAttrs  Attributes
}

// New creates new Kraph and returns it
func New(client kubernetes.Interface) (*Kraph, error) {
	return &Kraph{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(0.0, 0.0),
		nodeMap:                 make(map[types.UID]*Node),
		client:                  client,
	}, nil
}

// NewNode creates new kraph node and returns it
func (k *Kraph) NewNode(name string) *Node {
	return &Node{
		Node: k.WeightedUndirectedGraph.NewNode(),
		Name: name,
	}
}

// Nodes returns all kraph graph nodes
func (k *Kraph) Nodes() graph.Nodes {
	return k.WeightedUndirectedGraph.Nodes()
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

// Edges returns all kraph graph edges
func (k *Kraph) Edges() graph.Edges {
	return k.WeightedUndirectedGraph.Edges()
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

// Build builds the kubernetes resource graph
func (k *Kraph) Build() error {
	resources := []string{
		"nodes",
		"namespaces",
		"deployments",
		"replicasets",
		"daemonsets",
		"pods",
	}

	for _, r := range resources {
		if err := k.buildNodes(r); err != nil {
			return fmt.Errorf("failed building %s graph: %w", r, err)
		}
	}

	return nil
}

// buildNodes adds Kubernetes API objects to graph
func (k *Kraph) buildNodes(kind string) error {
	// simple options for now
	options := metav1.ListOptions{
		Limit: 100,
	}

	switch kind {
	case "nodes":
		n, err := k.client.CoreV1().Nodes().List(options)
		if err != nil {
			return fmt.Errorf("failed getting namespaces: %w", err)
		}
		return k.addNodes(n)
	case "namespaces":
		ns, err := k.client.CoreV1().Namespaces().List(options)
		if err != nil {
			return fmt.Errorf("failed getting namespaces: %w", err)
		}
		return k.addNamespaces(ns)
	case "deployments":
		dep, err := k.client.AppsV1().Deployments(metav1.NamespaceAll).List(options)
		if err != nil {
			return fmt.Errorf("failed getting replicasets: %w", err)
		}
		return k.addDeployments(dep)
	case "replicasets":
		rs, err := k.client.AppsV1().ReplicaSets(metav1.NamespaceAll).List(options)
		if err != nil {
			return fmt.Errorf("failed getting replicasets: %w", err)
		}
		return k.addReplicaSets(rs)
	case "daemonsets":
		ds, err := k.client.AppsV1().DaemonSets(metav1.NamespaceAll).List(options)
		if err != nil {
			return fmt.Errorf("failed getting replicasets: %w", err)
		}
		return k.addDaemonSets(ds)
	case "pods":
		p, err := k.client.CoreV1().Pods(metav1.NamespaceAll).List(options)
		if err != nil {
			return fmt.Errorf("failed getting pods: %v", err)
		}
		return k.addPods(p)
	default:
		return ErrUnknownObject
	}
}

// addEdges adds edges to the given node from all the owner nodes
func (k *Kraph) addEdges(to *Node, uid types.UID, owners []metav1.OwnerReference, weight float64) {
	for _, owner := range owners {
		//fmt.Println(to.Name, "Owners", owners)
		if from, ok := k.nodeMap[owner.UID]; ok {
			//fmt.Printf("Linking %s to %s", to.Name, from.Name)
			k.NewEdge(from, to, weight)
		}
	}
}

// linkNode links the node to its owners
func (k *Kraph) linkNode(name string, uid types.UID, owners []metav1.OwnerReference, weight float64) {
	if n, ok := k.nodeMap[uid]; ok {
		if kn := k.Node(n.ID()); kn == nil {
			node := k.NewNode(name)
			k.AddNode(node)
			k.addEdges(node, uid, owners, weight)
			return
		}
	}

	node := k.NewNode(name)
	k.AddNode(node)
	k.addEdges(node, uid, owners, weight)
	k.nodeMap[uid] = node
}

// addNodes gets a list of kubernetes nodes and adds them to kraph
func (k *Kraph) addNodes(nodes *corev1.NodeList) error {
	for _, node := range nodes.Items {
		k.linkNode(node.Name, node.UID, node.OwnerReferences, 0.0)
	}

	return nil
}

// addNamespaces gets a list of all kubernetes namespaces and adds them to kraph
func (k *Kraph) addNamespaces(namespaces *corev1.NamespaceList) error {
	for _, ns := range namespaces.Items {
		k.linkNode(ns.Name, ns.UID, ns.OwnerReferences, 0.0)
	}

	return nil
}

// addDeployments adds all kubernetes deployments to kraph
func (k *Kraph) addDeployments(dep *appsv1.DeploymentList) error {
	for _, d := range dep.Items {
		k.linkNode(d.Name, d.UID, d.OwnerReferences, 0.0)
	}

	return nil
}

// addReplicaSets adds all kubernetes replicasets to kraph
func (k *Kraph) addReplicaSets(rs *appsv1.ReplicaSetList) error {
	for _, r := range rs.Items {
		k.linkNode(r.Name, r.UID, r.OwnerReferences, 0.0)
	}

	return nil
}

// addDaemonSets adds all kubernetes DaemonSets to kraph
func (k *Kraph) addDaemonSets(ds *appsv1.DaemonSetList) error {
	for _, d := range ds.Items {
		k.linkNode(d.Name, d.UID, d.OwnerReferences, 0.0)
	}

	return nil
}

// addPods gets a list of all pods and adds them to kraph
func (k *Kraph) addPods(pods *corev1.PodList) error {
	for _, p := range pods.Items {
		k.linkNode(p.Name, p.UID, p.OwnerReferences, 0.0)
	}

	return nil
}
