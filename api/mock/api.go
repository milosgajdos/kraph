package mock

import (
	"github.com/milosgajdos/kraph/api/generic"
)

// API provides mock API
type API struct {
	*generic.API
}

// NewAPI returns new mock API
func NewAPI() *API {
	return &API{
		API: generic.NewAPI(),
	}
}
