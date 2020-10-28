package gh

import "github.com/milosgajdos/kraph/pkg/api/gen"

// API is GH API
type API struct {
	*gen.API
}

// NewAPI returns new K8s API object
func NewRepoAPI(s string) *API {
	return &API{
		API: gen.NewAPI(s),
	}
}
