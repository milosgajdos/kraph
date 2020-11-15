package types

// Resource is an API resource
type Resource struct {
	Name       string                 `json:"name"`
	Kind       string                 `json:"kind"`
	Group      string                 `json:"group"`
	Version    string                 `json:"version"`
	Namespaced bool                   `json:"namespaced"`
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
