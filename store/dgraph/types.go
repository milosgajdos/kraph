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
	Link      []Node    `json:"link,omitempty"`
	DType     []string  `json:"dgraph.type,omitempty"`

	// These are link facets
	Relation string  `json:"link|type,omitempty"`
	Weight   float64 `json:"link|weight,omitempty"`
}
