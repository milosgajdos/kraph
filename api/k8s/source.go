package k8s

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/gen"
)

// Source is API source
type Source struct {
	*gen.Source
}

// NewSource returns api.Source for k8s api
func NewSource(s string) api.Source {
	return &Source{
		Source: gen.NewSource(s),
	}
}
