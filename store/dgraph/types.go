package dgraph

// Node is dgraph node
type Node struct {
	UID  string `json:"xid"`
	Name string `json:"name"`
	Kind string `json:"kind"`
	Ns   string `json:"ns"`
}
