package star

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gen"
)

// Top is GitHub API starred repo topology
type Top struct {
	*gen.Top
}

// NewTop creates a new empty topology and returns it
func NewTop(a api.API) *Top {
	return &Top{
		Top: gen.NewTop(a),
	}
}
