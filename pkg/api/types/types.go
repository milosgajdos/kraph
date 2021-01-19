package types

// Resource is an API resource
type Resource struct {
	Group      string                 `json:"group"`
	Version    string                 `json:"version"`
	Kind       string                 `json:"kind"`
	Namespaced bool                   `json:"namespaced"`
	Name       string                 `json:"name"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// Link is an API link
type Link struct {
	UID      string                 `json:"uid"`
	From     string                 `json:"from"`
	To       string                 `json:"to"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Object is an API object
type Object struct {
	UID       string                 `json:"uid"`
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Resource  Resource               `json:"resource"`
	Links     []Link                 `json:"links"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
