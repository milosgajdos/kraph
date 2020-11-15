package star

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gen"
)

// API is a GitHub starred repo API
type API struct {
	*gen.API
}

// NewAPI creates a new GitHub starred repos API object and returns it
func NewRepoAPI(s api.Source) *API {
	return &API{
		API: gen.NewAPI(s),
	}
}
