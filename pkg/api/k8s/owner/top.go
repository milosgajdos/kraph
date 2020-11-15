package owner

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gen"
)

// Top is Kubernetes API topology
type Top struct {
	*gen.Top
}

// NewTop creates a new empty topology and returns it
func NewTop(a api.API) *Top {
	return &Top{
		Top: gen.NewTop(a),
	}
}
