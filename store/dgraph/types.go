package dgraph

import "time"

// Node is dgraph node
type Node struct {
	UID       string    `json:"uid,omitempty"`
	XID       string    `json:"xid,omitempty"`
	Name      string    `json:"name,omitempty"`
	Kind      string    `json:"kind,omitempty"`
	Namespace string    `json:"namespace,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	IsOwned   []Node    `json:"is_owned,omitempty"`
	DType     []string  `json:"dgraph.type,omitempty"`
}
