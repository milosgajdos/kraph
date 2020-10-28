package k8s

import "github.com/milosgajdos/kraph/pkg/api/gen"

// API is kubernetes API
type API struct {
	*gen.API
}

// NewAPI returns new K8s API object
func NewAPI(s string) *API {
	return &API{
		API: gen.NewAPI(s),
	}
}
