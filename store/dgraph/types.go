package dgraph

// Node is dgraph node
type Node struct {
	UID       string   `json:"xid"`
	Name      string   `json:"name"`
	Kind      string   `json:"kind"`
	Namespace string   `json:"namespace"`
	DType     []string `json:"dgraph.type,omitempty"`
}
