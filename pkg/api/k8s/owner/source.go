package owner

import "github.com/milosgajdos/kraph/pkg/api/generic"

// Source is the source of Github API.
type Source struct {
	*generic.Source
}

// NewSource returns API source.
func NewSource(s string) *Source {
	return &Source{
		Source: generic.NewSource(s),
	}
}
