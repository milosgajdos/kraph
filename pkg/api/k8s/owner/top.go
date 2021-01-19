package owner

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/generic"
)

// Top is Kubernetes API topology
type Top struct {
	*generic.Top
}

// NewTop creates a new empty topology and returns it
func NewTop(a api.API) *Top {
	return &Top{
		Top: generic.NewTop(a),
	}
}
