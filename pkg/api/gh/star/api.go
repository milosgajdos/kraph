package star

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/generic"
)

// API is a GitHub starred repo API
type API struct {
	*generic.API
}

// NewAPI creates a new GitHub starred repos API object and returns it
func NewRepoAPI(s api.Source) *API {
	return &API{
		API: generic.NewAPI(s),
	}
}
