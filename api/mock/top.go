package mock

import "github.com/milosgajdos/kraph/api/generic"

// Top provides mock Topology
type Top struct {
	*generic.Top
}

func NewTop() *Top {
	return &Top{
		Top: generic.NewTop(),
	}
}
