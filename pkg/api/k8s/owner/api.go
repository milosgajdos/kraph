package owner

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/generic"
)

// API is kubernetes API
type API struct {
	*generic.API
}

// NewAPI returns new K8s API object
func NewAPI(s api.Source) *API {
	return &API{
		API: generic.NewAPI(s),
	}
}
